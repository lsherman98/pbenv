package routes

import (
	"encoding/json"
	"time"

	"github.com/lsherman98/yt-rss/pocketbase/collections"
	"github.com/lsherman98/yt-rss/pocketbase/system"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
	"github.com/pocketbase/pocketbase/tools/types"
)

func renderStatsPageHandler(e *core.RequestEvent) error {
	html, err := template.NewRegistry().LoadFiles(
		"views/layout.html",
		"views/stats.html",
	).Render(nil)
	if err != nil {
		return e.BadRequestError("failed to load html", err)
	}
	return e.HTML(200, html)
}

func getStatsHandler(e *core.RequestEvent) error {
	authHeader := e.Request.Header.Get("Authorization")
	if authHeader == "" || !e.Auth.IsSuperuser() {
		return apis.NewApiError(403, "forbidden", nil)
	}

	stats, err := system.GetStats()
	if err != nil {
		return e.InternalServerError("failed to retrieve stats", err)
	}

	return e.JSON(200, stats)
}

type HistoricalCPU struct {
	Percent        float64        `json:"percent"`
	Created        types.DateTime `json:"created"`
	ProcessPercent float64        `json:"process_percent"`
}

type HistoricalMemory struct {
	Used           float64        `json:"used"`
	Total          float64        `json:"total"`
	Usage          float64        `json:"usage"`
	ProcessPercent float32        `json:"process_percent"`
	Created        types.DateTime `json:"created"`
}

type HistoricalDisk struct {
	Total   float64        `json:"total"`
	Used    float64        `json:"used"`
	Usage   float64        `json:"usage"`
	Created types.DateTime `json:"created"`
}

type HistoricalRuntime struct {
	Alloc      float64        `json:"alloc"`
	TotalAlloc float64        `json:"total_alloc"`
	Created    types.DateTime `json:"created"`
}

type HistoricalStats struct {
	CPU     []HistoricalCPU     `json:"cpu"`
	Memory  []HistoricalMemory  `json:"memory"`
	Disk    []HistoricalDisk    `json:"disk"`
	Runtime []HistoricalRuntime `json:"runtime"`
}

func getHistoricalStatsHandler(e *core.RequestEvent) error {
	authHeader := e.Request.Header.Get("Authorization")
	if authHeader == "" || !e.Auth.IsSuperuser() {
		return apis.NewApiError(403, "forbidden", nil)
	}

	period := e.Request.URL.Query().Get("period")
	var cutoff time.Time
	switch period {
	case "hour":
		cutoff = time.Now().Add(-1 * time.Hour)
	case "sixHrs":
		cutoff = time.Now().Add(-6 * time.Hour)
	case "day":
		cutoff = time.Now().Add(-24 * time.Hour)
	case "week":
		cutoff = time.Now().Add(-7 * 24 * time.Hour)
	case "fortnight":
		cutoff = time.Now().Add(-14 * 24 * time.Hour)
	}

	records, err := e.App.FindRecordsByFilter(collections.SystemStats, "created > {:cutoff}", "created", 0, 0, dbx.Params{
		"cutoff": cutoff.UTC(),
	})
	if err != nil {
		return e.InternalServerError("failed to retrieve historical stats", err)
	}

	e.App.Logger().Info("Retrieved historical stats", "period", period, "cutoff", cutoff, "records", len(records))

	cpuStats := make([]HistoricalCPU, 0, len(records))
	memoryStats := make([]HistoricalMemory, 0, len(records))
	diskStats := make([]HistoricalDisk, 0, len(records))
	runtimeStats := make([]HistoricalRuntime, 0, len(records))

	for _, record := range records {
		var stat system.SystemStats
		data := record.Get("data")
		jsonData, err := json.Marshal(data)
		if err != nil {
			return e.InternalServerError("failed to marshal record data", err)
		}

		if err := json.Unmarshal(jsonData, &stat); err != nil {
			return e.InternalServerError("failed to unmarshal system stats", err)
		}

		created := record.GetDateTime("created")

		cpuStats = append(cpuStats, HistoricalCPU{
			Percent:        stat.CPUPercent,
			ProcessPercent: stat.ProcessCPUPercent,
			Created:        created,
		})

		memoryStats = append(memoryStats, HistoricalMemory{
			Total:          float64(stat.Memory.Total),
			Usage:          stat.Memory.UsedPercent,
			Used:           float64(stat.Memory.Used),
			ProcessPercent: stat.ProcessMemoryPercent,
			Created:        created,
		})

		diskStats = append(diskStats, HistoricalDisk{
			Usage:   stat.Disk.UsedPercent,
			Total:   float64(stat.Disk.Total),
			Used:    float64(stat.Disk.Used),
			Created: created,
		})

		runtimeStats = append(runtimeStats, HistoricalRuntime{
			Alloc:      float64(stat.Runtime.Alloc),
			TotalAlloc: float64(stat.Runtime.TotalAlloc),
			Created:    created,
		})
	}

	data := &HistoricalStats{
		CPU:     cpuStats,
		Memory:  memoryStats,
		Disk:    diskStats,
		Runtime: runtimeStats,
	}

	return e.JSON(200, data)
}
