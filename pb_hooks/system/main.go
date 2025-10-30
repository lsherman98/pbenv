package system

import (
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

type SystemStats struct {
	CPU                   []cpu.InfoStat         `json:"cpu"`
	CPUPercent            float64                `json:"cpu_percent"`
	ProcessCPUPercent     float64                `json:"process_cpu_percent"`
	ProcessMemoryPercent  float32                `json:"process_memory_percent"`
	ProcessMemoryAbsolute float64                `json:"process_memory_absolute"`
	Memory                *mem.VirtualMemoryStat `json:"memory"`
	Swap                  *mem.SwapMemoryStat    `json:"swap"`
	Disk                  *disk.UsageStat        `json:"disk"`
	Host                  *host.InfoStat         `json:"host"`
	Runtime               *runtime.MemStats      `json:"runtime"`
}

func GetStats() (*SystemStats, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	var processCPUPercent float64
	var processMemoryPercent float32
	var processMemoryAbsolute float64

	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	processCPUPercent, err = p.CPUPercent()
	if err != nil {
		return nil, err
	}

	processMemoryPercent, err = p.MemoryPercent()
	if err != nil {
		return nil, err
	}

	memoryInfo, err := p.MemoryInfo()
	if err != nil {
		return nil, err
	}
	processMemoryAbsolute = float64(memoryInfo.RSS)

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	hostStat, err := host.Info()
	if err != nil {
		return nil, err
	}

	var runtimeStats runtime.MemStats
	runtime.ReadMemStats(&runtimeStats)

	return &SystemStats{
		CPU:                   cpuStat,
		CPUPercent:            cpuPercent[0],
		ProcessCPUPercent:     processCPUPercent,
		ProcessMemoryPercent:  processMemoryPercent,
		ProcessMemoryAbsolute: processMemoryAbsolute,
		Memory:                memStat,
		Swap:                  swapStat,
		Disk:                  diskStat,
		Host:                  hostStat,
		Runtime:               &runtimeStats,
	}, nil
}
