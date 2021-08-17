package main

import (
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	s "github.com/brenordv/azure-eventhub-tools/internal/senders"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	u.PrintHeader("SENDER")
	cfgFile := p.ParseCommandLine()
	u.LoadRuntimeConfig(cfgFile, p.ParseForSending)

	if d.CurrentConfig.OutboundConfig.Buffered {
		s.SendBuffered(d.CurrentConfig.OutboundConfig.OutboundFolder)
	} else {
		s.SendUnbuffered(d.CurrentConfig.OutboundConfig.OutboundFolder)
	}

	log.Printf("All done! (elapsed time: %s)\n", time.Since(start))
	os.Exit(d.ExitCode)
}
