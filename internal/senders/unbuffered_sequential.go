package senders

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	b "github.com/brenordv/azure-eventhub-tools/internal/builders"
	c "github.com/brenordv/azure-eventhub-tools/internal/clients"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"github.com/dgraph-io/badger/v3"
	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
	"log"
)

// TODO: add summary
func SendUnbuffered(t string) {
	bar := progressbar.Default(-1, "Sending unbuffered messages")
	defer h.CloseWithErrorHandling(bar.Close, "Failed to close progress bar.", false)
	db := c.OpenConnection(
		d.CurrentConfig.BadgerConfig.OutboundBaseDir,
		d.CurrentConfig.BadgerConfig.OutboundDir,
		d.CurrentConfig.BadgerConfig.OutboundValueDir)
	defer h.CloseWithErrorHandling(db.Close, "Failed to close database connection.", true)
	ctx, hub := c.GetEventHubClient(d.CurrentConfig.EventhubConnectionString, d.CurrentConfig.EntityPath)
	count := 0
	readQueue := make(chan bool, 100)
	sendAll := d.CurrentConfig.OutboundConfig.IgnoreStatus
	err := godirwalk.Walk(t, &godirwalk.Options{
		Unsorted: true,
		Callback: func(filename string, directoryEntry *godirwalk.Dirent) error {
			if !u.IsFile(filename) || (!sendAll && c.IsMessageProcessed(db, filename, true)) {
				h.DelegateIgnoreError(bar.Add, 0)
				return nil
			}

			readQueue <- true
			go func(f string, ch chan bool, bar *progressbar.ProgressBar, hub *eventhub.Hub,
				ctx context.Context, db *badger.DB) {
				defer h.DelegateIgnoreError(bar.Add, 1)
				defer u.FreeSlotOnQueue(ch)
				content := u.ReadTextFile(filename)
				msg := b.BuildOutboundMessage(filename, content)
				ev := eventhub.NewEventFromString(msg.Content)
				err := hub.Send(ctx, ev)

				if err == nil {
					msg.Status = d.Sent
				} else {
					msg.Status = d.Error
				}

				dbErr := db.Update(func(txn *badger.Txn) error {
					return txn.Set(msg.Id, p.SerializeOutboundMessage(&msg))
				})
				h.HandleError(fmt.Sprintf("Failed to buffer file '%s'.", f), dbErr, true)
			}(filename, readQueue, bar, hub, ctx, db)

			count++
			return nil
		},
		ErrorCallback: func(s string, err error) godirwalk.ErrorAction {
			h.HandleError(fmt.Sprintf("Failed to process file '%s'", s), err, true)
			return godirwalk.Halt
		},
	})

	for i := 0; i < cap(readQueue); i++ {
		readQueue <- true
	}
	log.Printf("Processed %d messages.", count)
	h.HandleError(fmt.Sprintf("Failed to walk through folder '%s'.", t), err, true)
}