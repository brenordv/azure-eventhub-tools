package readers

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	c "github.com/brenordv/azure-eventhub-tools/internal/clients"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	p "github.com/brenordv/azure-eventhub-tools/internal/parsers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"github.com/dgraph-io/badger/v3"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"strconv"
	"time"
)

var gdb *badger.DB
var pBar *progressbar.ProgressBar

func ReadFromEventHub() {
	if d.CurrentConfig.InboundConfig.ReadToFile {
		u.PrintReadAndSaveToDiskPerfWarning()
	}

	bar := progressbar.Default(-1, "Reading messages")
	defer h.CloseWithErrorHandling(bar.Close, "Failed to close progress bar.", false)
	db := c.OpenConnection(
		d.CurrentConfig.BadgerConfig.InboundBaseDir,
		d.CurrentConfig.BadgerConfig.InboundDir,
		d.CurrentConfig.BadgerConfig.InboundValueDir)
	defer h.CloseWithErrorHandling(db.Close, "Failed to close database connection.", true)

	gdb = db
	pBar = bar

	ctx, hub := c.GetEventHubClient(d.CurrentConfig.EventhubConnectionString, d.CurrentConfig.EntityPath)
	partitionId := strconv.Itoa(d.CurrentConfig.InboundConfig.PartitionId)

	if partitionId != "-1" {
		startReadingPartition(hub, ctx, partitionId, nil)
	} else {
		done := make(chan bool)
		for _, pId := range d.CurrentConfig.PartitionIds {
			go startReadingPartition(hub, ctx, pId, done)
		}
		<-done
	}
}

// OnMsgReceived is the handler for received messages on eventhub.
//
// Parameters:
//  _: Context. Passed automatically by the eventhub client. Not used, but can't get rid of it.
//  event: pointer to the event containing all the data we need.
//
// Returns:
//  Nothing
func OnMsgReceived(_ context.Context, event *eventhub.Event) error {
	now := time.Now()
	msg := d.InboundMessage{
		Id:             []byte(event.ID),
		EventId:        event.ID,
		QueuedTime:     *event.SystemProperties.EnqueuedTime,
		EventSeqNumber: event.SystemProperties.SequenceNumber,
		EventOffset:    event.SystemProperties.Offset,
		ProcessedAt:    now,
		MsgData:        string(event.Data),
	}
	defer h.DelegateIgnoreError(pBar.Add, 1)

	if shouldSaveMessage(msg.MsgData) {
		msg.SuggestedFilename = u.PutFileInSubFolderBasedOnTime(d.CurrentConfig.InboundConfig.InboundFolder,
			fmt.Sprintf("%s.txt", event.ID), now)

		err := DumpMessage(msg)
		if err != nil {
			return err
		}
		msg.Status = d.Exported
	}

	dbErr := gdb.Update(func(txn *badger.Txn) error {
		if !d.CurrentConfig.InboundConfig.IgnoreCheckpoint {
			_, err := txn.Get(msg.Id)
			if err != nil && err != badger.ErrKeyNotFound {
				return err
			}
			if err == badger.ErrKeyNotFound {
				return nil
			}
		}
		return txn.Set(msg.Id, p.SerializeInboundMessage(&msg))
	})

	return dbErr
}

// DumpMessage creates a text file with the Message received.
// Will panic in case of failure.
//
// Parameters:
//  m: Message with data extracted from the eventhub event.
//
// Returns:
//  error or nil if everything went well.
func DumpMessage(m d.InboundMessage) error {
	var content string
	if d.CurrentConfig.InboundConfig.DumpContentOnly {
		content = m.MsgData
	} else {
		content = m.ToString()
	}

	file, err := os.Create(m.SuggestedFilename)
	if err != nil {
		return nil
	}
	defer h.CloseWithErrorHandling(file.Close, fmt.Sprintf("Failed to close file '%s'", m.SuggestedFilename),
		true)

	_, err = io.WriteString(file, content)
	if err != nil {
		return nil
	}

	err = file.Sync()
	return err
}

func shouldSaveMessage(md string) bool {
	if !d.CurrentConfig.InboundConfig.ReadToFile {
		return false
	}

	if d.CurrentConfig.InboundConfig.DumpFilter == nil {
		return true
	}

	return d.CurrentConfig.InboundConfig.ContentHasFilterKeywords(md)
}

// startReadingPartition will start reading messages from eventhub for a specific partition.
//
// Parameters:
//  hub: eventhub.Hub client that will be used.
//  ctx: context in use.
//  partitionId: partition id that will be used for reading messages.
//  ch: channel used for controller code flow when reading from all partitions.
//
// Returns:
//  Nothing
func startReadingPartition(hub *eventhub.Hub, ctx context.Context, partitionId string, ch chan bool) {
	var listenerHandler *eventhub.ListenerHandle
	var err error
	listenerHandler, err = hub.Receive(ctx, partitionId, OnMsgReceived, eventhub.ReceiveWithConsumerGroup(d.CurrentConfig.InboundConfig.ConsumerGroup))
	h.HandleError("Failed to establish reading connection to EventHub.", err, true)

	h.WatchForUserInterruption(func() {
		if listenerHandler == nil {
			return
		}
		lastErr := listenerHandler.Err()
		h.HandleError("Failed to process received message.", lastErr, true)

		err := listenerHandler.Close(ctx)
		h.HandleError("Failed to close message listener.", err, true)

		err = hub.Close(ctx)
		h.HandleError("Failed to close eventhub client.", err, true)
		if ch != nil {
			ch <- true
		}
	})
}
