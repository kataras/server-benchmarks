package main

import (
	"errors"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type systemInfo struct {
	Processor string
	RAM       string
	OS        string

	Bombardier string // e..g bombardier version v1.2.4 windows/amd64 (bombardier --version)
	Go         string // e.g. go1.14 (go version)
	Dotnet     string // e.g. 3.1.102 (dotnet --version)
	Node       string // e.g. v13.10.1 (node --version)
}

func getSystemInfo() systemInfo {
	cpuStat, err := cpu.Info()
	catch(err)

	vmStat, err := mem.VirtualMemory()
	catch(err)

	hostStat, err := host.Info()
	catch(err)

	return systemInfo{
		Processor:  cpuStat[0].ModelName,
		RAM:        formatMemory(float64(vmStat.Total)),
		OS:         hostStat.Platform,
		Go:         getToolVer("go version", ""),
		Bombardier: getToolVer("bombardier --version", "v1.2.4"),
		Dotnet:     getToolVer("dotnet --version", ""),
		Node:       getToolVer("node --version", ""),
	}
}

func getToolVer(command string, def string) string {
	cmd := newCmd(command)
	out, err := cmd.CombinedOutput() // go: os.Stdout, bombardier: os.Stderr
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) { // when the tool is not installed, don't panic, just return an empty string.
			return "" // the template will omit this one. This may happen on the two optionals: dotnet and node.
		} else {
			v := string(out)
			if strings.HasPrefix(v, "The command could not be loaded") {
				// The new windows 11 on dotnet --version returns a specific error when dotnet was not installed.
				return ""
			}

			err = fmt.Errorf("command: %s\n%s\n%w", command, string(out), err)
			catch(err)
		}
	}

	v := strings.TrimPrefix(string(out), cmd.Args[0]+" version ")
	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, "\r", "")
	v = strings.TrimSpace(v)
	if spaceIdx := strings.IndexByte(v, ' '); spaceIdx > 0 {
		// e.g. go1.14 windows/amd64, remove the windows/amd64 (it's already on the OS info).
		v = v[:spaceIdx]
	}

	// If unspecified, return the default or the latest one, why?
	// Because, when go-get the bombardier tool, and not installed by its release page
	// then the version is not specified, this is probably a 3rd party bug-
	// I assume because the git commit info is missing.
	if v == "unspecified" {
		return def
	}

	return v
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func formatMemory(size float64) string {
	sizeSuffixes := [...]string{
		"B",
		"KB",
		"MB",
		"GB",
	}

	base := math.Log(size) / math.Log(1024)
	size = round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	suffix := sizeSuffixes[int(math.Floor(base))]
	return strconv.FormatFloat(size, 'f', -1, 64) + " " + string(suffix)
}
