package handlers

import (
	"os"
	"os/signal"
)

// WatchForUserInterruption will wait for the user to stop execution to continue.
// this means that any code after this method will only run after the ser presses Ctrl+C or something like that.
// Will panic in case of failure.
//
// Parameters:
//  None
//
// Returns:
//  Nothing.
func WatchForUserInterruption(cb func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	if cb == nil {
		return
	}
	cb()
}
