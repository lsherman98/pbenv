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


<img width="2918" height="1894" alt="127 0 0 1_8090___env" src="https://github.com/user-attachments/assets/b1d5aee0-bafc-48b0-b30a-2e62319c4fdc" />
<img width="3420" height="3054" alt="screencapture-127-0-0-1-8090-stats-2025-10-30-15_04_38" src="https://github.com/user-attachments/assets/82694a0f-f663-44f5-b068-eedfec9fe800" />
<img width="1710" height="948" alt="Screenshot 2025-10-31 at 1 00 02 AM" src="https://github.com/user-attachments/assets/7c5d3489-012d-4a1a-95d2-30ecb1fa6f18" />
