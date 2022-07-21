package senders

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	c "github.com/brenordv/azure-eventhub-tools/internal/clients"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	"github.com/dgraph-io/badger/v3"
	"github.com/schollz/progressbar/v3"
	"log"
	"sync"
)

// TODO: add summary
func SendBuffered(t string) {
	if !d.CurrentConfig.OutboundConfig.JustSendBuffered {
		count := BufferMessages(t)
		log.Printf("Total buffered files: %d\n", count)
	}
	sendMessagesSequentially()
}

// TODO: add summary
func sendMessagesSequentially() {
	bar := progressbar.Default(-1, "Sending buffered messages")
	defer h.CloseWithErrorHandling(bar.Close, "Failed to close progress bar.", false)
	db := c.OpenConnection(
		d.CurrentConfig.BadgerConfig.OutboundBaseDir,
		d.CurrentConfig.BadgerConfig.OutboundDir,
		d.CurrentConfig.BadgerConfig.OutboundValueDir)
	defer h.CloseWithErrorHandling(db.Close, "Failed to close database connection.", true)
	var wg sync.WaitGroup
	ctx, hub := c.GetEventHubClient(d.CurrentConfig.EventhubConnectionString, d.CurrentConfig.EntityPath)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()
		sendAll := d.CurrentConfig.OutboundConfig.IgnoreStatus
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				msg := p.DeserializeOutboundMessage(v)

				if !sendAll && msg.Status == d.Sent {
					h.DelegateIgnoreError(bar.Add, 0)
					return nil
				}

				go func(m *d.OutboundMessage, ctx context.Context, hub *eventhub.Hub, db *badger.DB,
					wg *sync.WaitGroup, bar *progressbar.ProgressBar, pKey *string) {
					wg.Add(1)
					defer wg.Done()
					defer h.DelegateIgnoreError(bar.Add, 1)
					ev := eventhub.NewEventFromString(m.Content)
					if pKey != nil {
						ev.PartitionKey = pKey
					}
					err := hub.Send(ctx, ev)

					dbErr := db.Update(func(txn *badger.Txn) error {
						if err == nil {
							m.Status = d.Sent
						} else {
							m.Status = d.Error
						}
						return txn.Set(m.Id, p.SerializeOutboundMessage(m))
					})
					h.HandleError(fmt.Sprintf("Failed to update message for file '%s'", m.FullFilename), dbErr, false)
					h.HandleError(fmt.Sprintf("Failed to send event for file '%s'", m.FullFilename), err, true)
				}(msg, ctx, hub, db, &wg, bar, d.CurrentConfig.OutboundConfig.PartitionKey)

				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	wg.Wait()
	h.HandleError("Failed to sequentially send buffered messages.", err, true)
}
