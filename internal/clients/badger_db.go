package clients

import (
	"fmt"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	"github.com/dgraph-io/badger/v3"
)

func OpenConnection(baseDir string, dir string, valueDir string) *badger.DB {
	var db *badger.DB
	var err error

	opts := badger.DefaultOptions(baseDir)
	opts.Dir = dir
	opts.ValueDir = valueDir
	opts.CompactL0OnClose = !d.CurrentConfig.BadgerConfig.BadgerSkipCompactL0OnClose
	opts.ValueLogFileSize = d.CurrentConfig.BadgerConfig.BadgerValueLogFileSize

	if !d.CurrentConfig.BadgerConfig.Verbose {
		opts.Logger = nil
	}

	db, err = badger.Open(opts)
	h.HandleError("Failed to open badger database.", err, true)
	return db
}

func ReadOutboundMessageFromDb(sTicket *badger.Item) (d.OutboundMessage, error) {
	var msg d.OutboundMessage
	err := sTicket.Value(func(val []byte) error {
		msg = *p.DeserializeOutboundMessage(val)
		return nil
	})
	return msg, err
}

func IsMessageProcessed(db *badger.DB, f string, allowErrorRetry bool) bool {
	var isProcessed bool
	dbErr := db.View(func(txn *badger.Txn) error {
		r, err := txn.Get([]byte(f))
		if err == badger.ErrKeyNotFound {
			return nil
		}

		msg, er := ReadOutboundMessageFromDb(r)
		h.HandleError(fmt.Sprintf("Failed to deserialize %s", f), er, true)

		if msg.Status != d.Sent || (msg.Status == d.Error && allowErrorRetry) {
			return nil
		}

		isProcessed = true
		return nil
	})

	h.HandleError(fmt.Sprintf("Failed to get message on db for file '%s'.\n", f),dbErr, true)
	return isProcessed
}
