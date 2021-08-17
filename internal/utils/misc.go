package utils

import (
	"fmt"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	"log"
)

// TODO: add summary
func PrintHeader(toolName string) {
	fmt.Printf("%sAzure Eventhub Tools%s: %s%s%s \n", d.ColorBlue, d.ColorReset, d.ColorYellow, toolName, d.ColorReset)
	fmt.Printf("(v: %s)\n", d.Version)
}

// PrintReadAndSafeToDiskPerfWarning simply prints out a warning message if the user decides
// to read messages AND write them to disk at the same time.
//
// Parameters:
//  None.
//
// Returns:
//  Nothing.
func PrintReadAndSafeToDiskPerfWarning() {
	log.Println("----| WARNING | ----------------------------------")
	log.Println("Reading to file drastically SLOWS things down.")
	log.Println("--------------------------------------------------")
}