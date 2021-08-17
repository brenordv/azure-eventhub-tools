package handlers

import (
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	"log"
	"runtime"
)

// HandleError just a wrapper to handle errors in a neat, lazy manner.
//
// Parameters:
//  customMsg: this message will be logger alongside the error
//  err: error object that will be logger.
//  shouldPanic: if true, will panic. Otherwise, will just signal the application to exit.
//
// Returns:
//  Nothing
func HandleError(customMsg string, err error, shouldPanic bool) {
	if err == nil {
		return
	}

	log.Printf("[ERROR] %s. Details: %s\n", customMsg, err)

	if shouldPanic {
		panic(err)
	}

	log.Println(err)
	d.ExitCode = 0
	runtime.Goexit()
}

// TODO: add summary
func CloseWithErrorHandling(fn d.FuncReturnError, msg string, shouldPanic bool) {
	err := fn()
	HandleError(msg, err, shouldPanic)
}

// TODO: add summary
func DelegateIgnoreError(fn d.FuncReturnErrorIntArg, n int) {
	_ = fn(n)
}