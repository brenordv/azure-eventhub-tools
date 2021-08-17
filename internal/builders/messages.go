package builders

import (
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	"time"
)

// TODO: add summary
func BuildOutboundMessage(f string, c string) d.OutboundMessage {
	return d.OutboundMessage{
		Id:           []byte(f),
		FullFilename: f,
		Content:      c,
		ProcessedAt:  time.Now(),
		Status:      d.New,
	}
}
