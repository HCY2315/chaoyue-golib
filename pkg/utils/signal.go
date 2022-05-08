package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitQuitSignals() os.Signal {
	return WaitSignals(syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
}

func WaitSignals(sigList ...os.Signal) os.Signal {
	if len(sigList) == 0 {
		panic("wait at least one signal")
	}
	// catch system signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, sigList...)
	return <-signals
}
