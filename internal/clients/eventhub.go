package clients

import (
	"context"
	"fmt"
	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	"log"
	"strings"
	"time"
)

// parseConnectionString parses connection string to make sure it will have the correct entityPath in it.
// Will panic in case of failure.
//
// Parameters:
//  connectionString: connection string that will be used to open a connection to Eventhub
//  entityPath: Name of the entity path (eventhub) that will be targeted.
//
// Returns:
//  Nothing. Use the global variables.
func parseConnectionString(connectionString string, entityPath string) string {
	if !strings.Contains(connectionString, ";EntityPath=") {
		return fmt.Sprintf("%s;EntityPath=%s", connectionString, entityPath)
	}
	return connectionString
}

// logRuntimeInfo will log the runtime information for each available partition.
// Will panic in case of failure.
//
// Parameters:
//  hub: pointer to the eventhub.Hub object.
//  ctx: Current context
//
// Returns:
//  Nothing. Use the global variables.
func logRuntimeInfo(hub *eventhub.Hub, ctx context.Context) {
	info, e := hub.GetRuntimeInformation(ctx)
	h.HandleError("Failed to get runtime information", e, false)

	log.Printf("Runtime started at '%s', pointing at path '%s' with %d partitions. Available partitions: %s\n",
		info.CreatedAt, info.Path, info.PartitionCount, info.PartitionIDs)
	for _, p := range info.PartitionIDs {
		pInfo, e := hub.GetPartitionInformation(ctx, p)
		h.HandleError(fmt.Sprintf("Failed to get info for partition '%s'.", p), e, true)
		log.Printf(
			"Partition: '%s'. HubPath: %s | BeginningSequenceNumber: %d | LastSequenceNumber: %d | LastEnqueuedOffset: %s | LastEnqueuedTimeUtc: %s\n",
			p,
			pInfo.HubPath,
			pInfo.BeginningSequenceNumber,
			pInfo.LastSequenceNumber,
			pInfo.LastEnqueuedOffset,
			pInfo.LastEnqueuedTimeUtc.Format(time.RFC3339Nano))
	}
}

// GetEventHubClient instantiate an Eventhub.
// Will panic in case of failure.
//
// Parameters:
//  connectionString: connection string that will be used to open a connection to Eventhub
//  entityPath: Name of the entity path (eventhub) that will be targeted.
//  renew: if true, will close existing connection and open a new one.
//
// Returns:
//  Nothing. Use the global variables.
func GetEventHubClient(connectionString string, entityPath string) (context.Context, *eventhub.Hub) {
	cs := parseConnectionString(connectionString, entityPath)
	ctx := context.Background()
	hub, err := eventhub.NewHubFromConnectionString(cs)
	h.HandleError("Failed to create new EventHub client from connection string", err, true)

	if !d.CurrentConfig.SkipGetRuntimeInfo {
		logRuntimeInfo(hub, ctx)
	}

	return ctx, hub
}
