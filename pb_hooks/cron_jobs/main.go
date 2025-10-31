package cron_jobs

import (
	"encoding/json"
	"time"

	"github.com/lsherman98/pbenv/pb_hooks/system"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Init(app *pocketbase.PocketBase) error {
	app.Cron().MustAdd("MeasureSystemStats", "* * * * *", func() {
		stats, err := system.GetStats()
		if err != nil {
			app.Logger().Error("Failed to get system stats:", "error", err)
			return
		}

		statsJson, err := json.Marshal(stats)
		if err != nil {
			return
		}

		systemStatsCollection, err := app.FindCollectionByNameOrId("system_stats")
		if err != nil {
			return
		}

		statsRecord := core.NewRecord(systemStatsCollection)
		statsRecord.Set("data", string(statsJson))
		if err := app.Save(statsRecord); err != nil {
			return
		}
	})

	app.Cron().MustAdd("CleanUpSystemStats", "0 0 * * *", func() {
		systemStatsCollection, err := app.FindCollectionByNameOrId("system_stats")
		if err != nil {
			return
		}

		cutoffDate := time.Now().AddDate(0, 0, -14).UTC()
		records, err := app.FindRecordsByFilter(systemStatsCollection, "created < {:cutoff}", "", 0, 0, dbx.Params{
			"cutoff": cutoffDate.Format(time.RFC3339),
		})
		if err != nil {
			return
		}

		for _, record := range records {
			err := app.Delete(record)
			if err != nil {
				continue
			}
		}
	})

	return nil
}
