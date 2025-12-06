package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	robotgo.KeySleep = 50
	mux := http.NewServeMux()
	mux.HandleFunc("/start-sequence", func(w http.ResponseWriter, r *http.Request) {
		// Export current state by timestamp
		sname := "auto_" + time.Now().Format("2006-01-02_15-04-05")
		robotgo.KeyTap("alt", "f")
		for range 2 {
			robotgo.KeyTap("down")
		}
		robotgo.KeyTap("enter")
		robotgo.KeyTap("tab")
		for range 8 {
			robotgo.KeyTap("down")
		}
		robotgo.KeyTap("space")
		robotgo.KeyTap("down")
		robotgo.KeyTap("enter")
		robotgo.Sleep(1)
		robotgo.TypeStr(sname)
		robotgo.KeyTap("enter")
		robotgo.Sleep(1)

		// Start Match Timer
		robotgo.KeyTap("alt", "f")
		robotgo.KeyTap("escape")
		for range 9 {
			robotgo.KeyTap("tab")
		}
		robotgo.KeyTap("space")
		slog.Info("State saved and match started", "save", sname + ".bsm")
	})

	s := http.Server{
		Addr:    ":9269",
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error with webserver", "error", err)
			quit <- syscall.SIGINT

		}
	}()
	slog.Info("Startup Complete!")

	<-quit
	slog.Info("Shutting Down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Error during shutdown", "error", err)
		os.Exit(2)
	}
	slog.Info("Goodbye!")
}
