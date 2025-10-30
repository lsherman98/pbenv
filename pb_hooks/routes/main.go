package routes

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func Init(app *pocketbase.PocketBase) error {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		env := se.Router.Group("/_/env")
		env.GET("", renderEnvPageHandler)
		env.POST("", addEnvHandler).Bind(apis.RequireSuperuserAuth())
		env.GET("/data", getEnvsHandler).Bind(apis.RequireSuperuserAuth())
		env.PUT("/{key}", updateEnvHandler).Bind(apis.RequireSuperuserAuth())
		env.DELETE("/{key}", deleteEnvHandler).Bind(apis.RequireSuperuserAuth())

		stats := se.Router.Group("/_/stats")
		stats.GET("", renderStatsPageHandler)
		stats.GET("/data", getStatsHandler).Bind(apis.RequireSuperuserAuth())
		stats.GET("/historical", getHistoricalStatsHandler).Bind(apis.RequireSuperuserAuth())

		se.Router.POST("/_/restart", restartHandler).Bind(apis.RequireSuperuserAuth())

		return se.Next()
	})

	return nil
}
