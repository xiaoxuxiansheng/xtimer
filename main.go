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

	migratorApp.Start()
	schedulerApp.Start()
	defer schedulerApp.Stop()

	webServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
}
