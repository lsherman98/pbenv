package system

import (
	"os"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

var (
	prevNetworkStats  *net.IOCountersStat
	networkStatsMutex sync.Mutex
)

type NetworkStats struct {
	BytesRecv   uint64 `json:"bytes_received"`
	BytesSent   uint64 `json:"bytes_sent"`
	PacketsRecv uint64 `json:"packets_received"`
	PacketsSent uint64 `json:"packets_sent"`
}

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
	Runtime               *MemStats              `json:"runtime"`
	NetworkStats          NetworkStats           `json:"network_stats"`
}

type MemStats struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
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

	vmemStat, err := mem.VirtualMemory()
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

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	runtimeStats := MemStats{
		Alloc:      memStats.Alloc,
		TotalAlloc: memStats.TotalAlloc,
		Sys:        memStats.Sys,
		NumGC:      memStats.NumGC,
	}

	counters, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	var bytesSent, bytesRecv, packetsSent, packetsRecv uint64
	if len(counters) > 0 {
		currentStats := counters[0]

		networkStatsMutex.Lock()
		if prevNetworkStats != nil {
			if currentStats.BytesSent >= prevNetworkStats.BytesSent {
				bytesSent = currentStats.BytesSent - prevNetworkStats.BytesSent
			}
			if currentStats.BytesRecv >= prevNetworkStats.BytesRecv {
				bytesRecv = currentStats.BytesRecv - prevNetworkStats.BytesRecv
			}
			if currentStats.PacketsSent >= prevNetworkStats.PacketsSent {
				packetsSent = currentStats.PacketsSent - prevNetworkStats.PacketsSent
			}
			if currentStats.PacketsRecv >= prevNetworkStats.PacketsRecv {
				packetsRecv = currentStats.PacketsRecv - prevNetworkStats.PacketsRecv
			}
		}
		prevNetworkStats = &currentStats
		networkStatsMutex.Unlock()
	}

	return &SystemStats{
		CPU:                   cpuStat,
		CPUPercent:            cpuPercent[0],
		ProcessCPUPercent:     processCPUPercent,
		ProcessMemoryPercent:  processMemoryPercent,
		ProcessMemoryAbsolute: processMemoryAbsolute,
		Memory:                vmemStat,
		Swap:                  swapStat,
		Disk:                  diskStat,
		Host:                  hostStat,
		Runtime:               &runtimeStats,
		NetworkStats: NetworkStats{
			BytesRecv:   bytesRecv,
			BytesSent:   bytesSent,
			PacketsRecv: packetsRecv,
			PacketsSent: packetsSent,
		},
	}, nil
}
