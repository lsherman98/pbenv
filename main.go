package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lsherman98/pbenv/pb_hooks/cron_jobs"
	"github.com/lsherman98/pbenv/pb_hooks/routes"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := routes.Init(app); err != nil {
		log.Fatal(err)
	}

	if err := cron_jobs.Init(app); err != nil {
		log.Fatal(err)
	}

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
