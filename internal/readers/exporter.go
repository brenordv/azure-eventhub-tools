package readers

import (
	"fmt"
	c "github.com/brenordv/azure-eventhub-tools/internal/clients"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"github.com/dgraph-io/badger/v3"
	"github.com/schollz/progressbar/v3"
)

func ExportFromCheckpoint() {
	bar := progressbar.Default(-1, "Exporting messages")
	defer h.CloseWithErrorHandling(bar.Close, "Failed to close progress bar.", false)
	db := c.OpenConnection(
		d.CurrentConfig.BadgerConfig.InboundBaseDir,
		d.CurrentConfig.BadgerConfig.InboundDir,
		d.CurrentConfig.BadgerConfig.InboundValueDir)
	defer h.CloseWithErrorHandling(db.Close, "Failed to close database connection.", true)
	readQueue := make(chan bool, 100)

	dbErr := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()


		for it.Rewind(); it.Valid(); it.Next() {
			var msg d.InboundMessage
			item := it.Item()
			_ = item.Value(func(v []byte) error {
				msg = *p.DeserializeInboundMessage(v)
				return nil
			})

			msg.SuggestedFilename= u.PutFileInSubFolderBasedOnTime(d.CurrentConfig.InboundConfig.InboundFolder,
				fmt.Sprintf("%s.txt", msg.EventId), msg.ProcessedAt)

			if u.Exists(msg.SuggestedFilename) {
				h.DelegateIgnoreError(bar.Add, 0)
				continue
			}

			err := DumpMessage(msg)
			if err != nil {
				return err
			}
			msg.Status = d.Exported

			err = txn.Set(msg.Id, p.SerializeInboundMessage(&msg))
			if err == badger.ErrTxnTooBig {
				err = txn.Commit()
				h.HandleError("Failed to commit transaction.", err, true)

				err = txn.Set(msg.Id, p.SerializeInboundMessage(&msg))
				h.HandleError("Failed to retry update inbound message after committing.", err, true)
			}
			h.DelegateIgnoreError(bar.Add, 1)
		}
		return nil
	})

	for i := 0; i < cap(readQueue); i++ {
		readQueue <- true
	}
	h.HandleError("Failed to export messages", dbErr, true)
}
