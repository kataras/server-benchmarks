package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
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
		Go:         getToolVer("go version"),
		Bombardier: getToolVer("bombardier --version"),
		Dotnet:     getToolVer("dotnet --version"),
		Node:       getToolVer("node --version"),
	}
}

func getToolVer(command string) string {
	cmd := newCmd(command)
	out, err := cmd.CombinedOutput() // go: os.Stdout, bombardier: os.Stderr
	catch(err)
	v := strings.TrimPrefix(string(out), cmd.Args[0]+" version ")
	v = strings.ReplaceAll(v, "\n", "")
	v = strings.ReplaceAll(v, "\r", "")
	v = strings.TrimSpace(v)
	if spaceIdx := strings.IndexByte(v, ' '); spaceIdx > 0 {
		// e.g. go1.14 windows/amd64, remove the windows/amd64 (it's already on the OS info).
		v = v[:spaceIdx]
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
