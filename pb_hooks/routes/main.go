package routes

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Init(app *pocketbase.PocketBase) error {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		env := se.Router.Group("/_/env")
		env.GET("", renderEnvPageHandler)
		env.POST("", addEnvHandler)
		env.GET("/data", getEnvsHandler)
		env.PUT("/{key}", updateEnvHandler)
		env.DELETE("/{key}", deleteEnvHandler)

		stats := se.Router.Group("/_/stats")
		stats.GET("", renderStatsPageHandler)
		stats.GET("/data", getStatsHandler)
		stats.GET("/historical", getHistoricalStatsHandler)
		
		se.Router.POST("/_/restart", restartHandler)
		
		return se.Next()
	})

	return nil
}
