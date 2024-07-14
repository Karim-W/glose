package glose

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Closable interface {
	Close() error
}

func Panik(err error) {
	if err == nil {
		return
	}

	slog.Error("Panic: ", slog.String("error", err.Error()))

	for _, closable := range ClosableList {
		if err := closable.Close(); err != nil {
			slog.Error("Error closing: ", slog.String("error", err.Error()))
		}
	}

	panic(err)
}

var ClosableList []Closable = make([]Closable, 0)

func Watch(
	closables ...Closable,
) {
	ClosableList = append(ClosableList, closables...)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down...")

	for _, closable := range ClosableList {
		if err := closable.Close(); err != nil {
			slog.Error("Error closing: ", slog.String("error", err.Error()))
		}
	}

	slog.Info("Shutdown complete")

	if os.Getenv("GLOSE_SKIP_EXIT") != "" {
		return
	}
	os.Exit(0)
}

func Register(closables ...Closable) {
	ClosableList = append(ClosableList, closables...)
}
