# PBENV

This is a Pocketbase plugin for managing your environment variables on your VPS. It also tracks system and process information so you can see current and historical resource consumption.

Requires extending pocketbase with Go.

## To start the example project

1. `git clone github.com/lsherman98/pbenv`
2. `go mod tidy`
3. To start the server: `make serve`

Visit `localhost:8090/_/stats` and `localhost:8090/_/env`

The server restart hook expects two environment variables:

-   `DEV` should be `true` or `false`
-   `RESTART_CMD` for exampke `systemctl restart pocketbase.service`

When `DEV` is `true`, the restart hook just kills the current process.

## To add this plugin to your Pocketbase project

1. Copy the `views` directory to the root of your project.
2. Copy `cron_jobs`, `routes`, and `system` directories to your `pb_hooks` folder.
3. Import the `system_stats` collection to your project; copy the JSON from `collections.json`.

### If you don't want to persist historical system metrics:

-   No need to import the new collection.
-   Delete the cron jobs.
-   Update the JavaScript: remove the charts and JavaScript in `views/stats` that handles getting historical data.
-   Remove the script tags in `views/layout.html` that import chart.js
-   Delete the `GET /historical` route in `pb_hooks/routes/main.go`and associated handler.

## Credits

Inspired by:
- [pb-ext](https://github.com/magooney-loon/pb-ext)
- [pb-hooks-dash](https://github.com/deselected/pb-hooks-dash)