package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xiaoxuxiansheng/xtimer/app"
)

func main() {
	migratorApp := app.GetMigratorApp()
	schedulerApp := app.GetSchedulerApp()
	webServer := app.GetWebServer()
	monitor := app.GetMonitorApp()

	migratorApp.Start()
	schedulerApp.Start()
	defer schedulerApp.Stop()

	monitor.Start()
	webServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
}
