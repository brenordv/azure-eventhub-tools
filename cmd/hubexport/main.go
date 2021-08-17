package main

import (
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	r "github.com/brenordv/azure-eventhub-tools/internal/readers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	u.PrintHeader("EXPORT")
	cfgFile := p.ParseCommandLine()
	u.LoadRuntimeConfig(cfgFile, p.ParseForExport)

	r.ExportFromCheckpoint()

	log.Printf("All done! (elapsed time: %s)\n", time.Since(start))
	os.Exit(d.ExitCode)
}
