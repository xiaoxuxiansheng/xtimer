package main

import (
	"net/http"
	_ "net/http/pprof"
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

	// 支持 pprof
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
		_ = http.ListenAndServe(":9999", nil)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
}
