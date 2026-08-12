package report

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/runner"
)

// The charts are self-contained SVG documents rendered without any chart
// library, so the output is deterministic, dependency-free and works
// wherever an <img> works (GitHub READMEs included). Colors follow the
// environment's *language* (a stable categorical assignment), values are
// printed at every bar tip — the chart is a static image, so labels do the
// job a tooltip would — and a prefers-color-scheme style block adapts the
// chart to light and dark pages.

// Chart geometry (pixels).
const (
	chartWidth   = 920
	chartPadX    = 20
	chartHeaderH = 78 // title + subtitle + legend.
	chartRowH    = 30
	chartBarH    = 18 // mark spec: bars stay under 24px thick.
	chartPadBot  = 18
	chartNameW   = 180 // right-aligned name column.
	chartValueW  = 100 // reserved space for the value labels at bar tips.
)

// chartFont keeps every chart in the system UI sans.
const chartFont = "system-ui, -apple-system, 'Segoe UI', sans-serif"

// chartStyle carries the light palette with a dark override: categorical
// series slots (validated for contrast and color-vision-deficiency
// separation on both surfaces), plus text and chrome tokens. Text never
// wears a series color.
const chartStyle = `text{font-family:` + chartFont + `}
.surface{fill:#fcfcfb;stroke:rgba(11,11,11,0.10)}
.title{fill:#0b0b0b;font-size:15px;font-weight:600}
.subtitle{fill:#898781;font-size:12px}
.name{fill:#0b0b0b;font-size:12.5px}
.value{fill:#52514e;font-size:12px}
.legend{fill:#52514e;font-size:11.5px}
.baseline{stroke:#c3c2b7;stroke-width:1}
.s0{fill:#2a78d6}.s1{fill:#eb6834}.s2{fill:#1baf7a}.s3{fill:#eda100}
.s4{fill:#e87ba4}.s5{fill:#008300}.s6{fill:#4a3aa7}.s7{fill:#e34948}
@media (prefers-color-scheme: dark){
.surface{fill:#1a1a19;stroke:rgba(255,255,255,0.10)}
.title{fill:#ffffff}
.name{fill:#ffffff}
.value{fill:#c3c2b7}
.legend{fill:#c3c2b7}
.baseline{stroke:#383835}
.s0{fill:#3987e5}.s1{fill:#d95926}.s2{fill:#199e70}.s3{fill:#c98500}
.s4{fill:#d55181}.s5{fill:#008300}.s6{fill:#9085e9}.s7{fill:#e66767}
}`

// chartBar is one bar of a chart.
type chartBar struct {
	name  string
	slot  int     // categorical color slot of the env's language.
	value float64 // drives the bar length; must be >= 0.
	label string  // printed at the bar tip.
}

// legendItem maps a language to its color slot.
type legendItem struct {
	label string
	slot  int
}

// languageSlots assigns each language a stable categorical color slot in
// order of first appearance across the whole report, so a language keeps
// its color in every chart of a run.
func languageSlots(tests []runner.TestReport) map[string]int {
	slots := make(map[string]int)
	for _, tr := range tests {
		for _, er := range tr.Envs {
			if _, ok := slots[er.Env.Language]; !ok {
				slots[er.Env.Language] = len(slots) % 8
			}
		}
	}

	return slots
}

// renderRPSChart renders the requests-per-second chart of a test.
// results must be the successful envs in report order (best first).
func renderRPSChart(testName string, results []runner.EnvReport, slots map[string]int) []byte {
	bars := make([]chartBar, 0, len(results))
	for _, er := range results {
		bars = append(bars, chartBar{
			name:  er.Env.Name,
			slot:  slots[er.Env.Language],
			value: er.Result.RequestsPerSecond.Mean,
			label: thousands(er.Result.RequestsPerSecond.Mean),
		})
	}

	return renderBarChart(testName, "Requests per second (higher is better)", bars, legendFor(results, slots))
}

// renderLatencyChart renders the mean-latency chart of a test, ranked
// best (lowest) first.
func renderLatencyChart(testName string, results []runner.EnvReport, slots map[string]int) []byte {
	ranked := slices.Clone(results)
	slices.SortStableFunc(ranked, func(a, b runner.EnvReport) int {
		if d := a.Result.Latency.Mean - b.Result.Latency.Mean; d != 0 {
			if d < 0 {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Env.Name, b.Env.Name)
	})

	bars := make([]chartBar, 0, len(ranked))
	for _, er := range ranked {
		bars = append(bars, chartBar{
			name:  er.Env.Name,
			slot:  slots[er.Env.Language],
			value: er.Result.Latency.Mean,
			label: bombardier.FormatTimeUs(er.Result.Latency.Mean),
		})
	}

	return renderBarChart(testName, "Mean latency (lower is better)", bars, legendFor(ranked, slots))
}

// legendFor lists the languages present in results, in slot order.
func legendFor(results []runner.EnvReport, slots map[string]int) []legendItem {
	seen := make(map[string]bool)
	var items []legendItem
	for _, er := range results {
		if lang := er.Env.Language; !seen[lang] {
			seen[lang] = true
			items = append(items, legendItem{label: lang, slot: slots[lang]})
		}
	}

	slices.SortFunc(items, func(a, b legendItem) int { return a.slot - b.slot })
	return items
}

// renderBarChart renders a horizontal bar chart as a standalone SVG
// document.
func renderBarChart(title, subtitle string, bars []chartBar, legend []legendItem) []byte {
	height := chartHeaderH + len(bars)*chartRowH + chartPadBot
	plotX := chartNameW
	plotWidth := chartWidth - plotX - chartValueW - chartPadX

	maxValue := 0.0
	for _, b := range bars {
		maxValue = max(maxValue, b.value)
	}

	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-label="%s">`,
		chartWidth, height, chartWidth, height, xmlEscape(title+": "+subtitle))
	svg.WriteString("\n")
	fmt.Fprintf(&svg, "<title>%s: %s</title>\n", xmlEscape(title), xmlEscape(subtitle))
	fmt.Fprintf(&svg, "<style>%s</style>\n", chartStyle)

	// Card surface.
	fmt.Fprintf(&svg, `<rect class="surface" x="0.5" y="0.5" width="%d" height="%d" rx="8"/>`+"\n", chartWidth-1, height-1)

	// Header: title, subtitle, legend at the top right.
	fmt.Fprintf(&svg, `<text class="title" x="%d" y="30">%s</text>`+"\n", chartPadX, xmlEscape(title))
	fmt.Fprintf(&svg, `<text class="subtitle" x="%d" y="50">%s</text>`+"\n", chartPadX, xmlEscape(subtitle))

	x := float64(chartWidth - chartPadX)
	for _, item := range slices.Backward(legend) {
		labelWidth := textWidth(item.label, 11.5)
		x -= labelWidth
		fmt.Fprintf(&svg, `<text class="legend" x="%.1f" y="30">%s</text>`+"\n", x, xmlEscape(item.label))
		x -= 16 // swatch plus gap.
		fmt.Fprintf(&svg, `<rect class="s%d" x="%.1f" y="21" width="10" height="10" rx="3"/>`+"\n", item.slot, x)
		x -= 18 // gap between legend entries.
	}

	// Bars.
	for i, b := range bars {
		rowY := float64(chartHeaderH + i*chartRowH)
		barY := rowY + float64(chartRowH-chartBarH)/2
		textY := barY + float64(chartBarH)/2 + 4.25

		barWidth := 2.0
		if maxValue > 0 {
			barWidth = max(2, b.value/maxValue*float64(plotWidth))
		}

		fmt.Fprintf(&svg, `<text class="name" text-anchor="end" x="%d" y="%.2f">%s</text>`+"\n",
			plotX-10, textY, xmlEscape(b.name))
		fmt.Fprintf(&svg, `<path class="s%d" d="%s"/>`+"\n", b.slot, barPath(float64(plotX), barY, barWidth, chartBarH))
		fmt.Fprintf(&svg, `<text class="value" x="%.1f" y="%.2f">%s</text>`+"\n",
			float64(plotX)+barWidth+8, textY, xmlEscape(b.label))
	}

	// Baseline the bars grow from.
	fmt.Fprintf(&svg, `<line class="baseline" x1="%.1f" y1="%d" x2="%.1f" y2="%d"/>`+"\n",
		float64(plotX)+0.5, chartHeaderH-6, float64(plotX)+0.5, chartHeaderH+len(bars)*chartRowH+2)

	svg.WriteString("</svg>\n")
	return []byte(svg.String())
}

// barPath draws a horizontal bar with a 4px-rounded data end and a square
// baseline end, growing right from (x, y).
func barPath(x, y, width, height float64) string {
	r := min(4.0, width/2)

	return fmt.Sprintf("M%.1f,%.1f h%.1f a%.1f,%.1f 0 0 1 %.1f,%.1f v%.1f a%.1f,%.1f 0 0 1 %.1f,%.1f h%.1f z",
		x, y,
		width-r,
		r, r, r, r,
		height-2*r,
		r, r, -r, r,
		-(width - r))
}

// textWidth approximates the rendered width of s in the chart font at the
// given size — enough precision to lay out legend entries and reserve
// label space.
func textWidth(s string, size float64) float64 {
	return float64(len(s)) * size * 0.58
}

// thousands formats a value as a grouped integer, e.g. 284059.44 -> "284,059".
func thousands(v float64) string {
	s := strconv.FormatFloat(v, 'f', 0, 64)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}

	return s
}

// xmlEscape escapes s for use in SVG text content and attributes.
func xmlEscape(s string) string {
	return xmlReplacer.Replace(s)
}

var xmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)
