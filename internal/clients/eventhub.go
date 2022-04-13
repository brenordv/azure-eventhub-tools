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

// extractRuntimeInfo will extract the partition ids and log the runtime information for each available partition, if
// parameter alsoLog is true.
// Will panic in case of failure.
//
// Parameters:
//  hub: pointer to the eventhub.Hub object.
//  ctx: Current context
//  alsoLog: if true will fetch and log details about each partition.
//
// Returns:
//  Slice of strings containing the PartitionIds available.
func extractRuntimeInfo(hub *eventhub.Hub, ctx context.Context, alsoLog bool) []string {
	info, e := hub.GetRuntimeInformation(ctx)
	h.HandleError("Failed to get runtime information", e, false)

	log.Printf("Runtime started at '%s', pointing at path '%s' with %d partitions. Available partitions: %s\n",
		info.CreatedAt, info.Path, info.PartitionCount, info.PartitionIDs)

	if !alsoLog {
		return info.PartitionIDs
	}

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

	return info.PartitionIDs
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

	d.CurrentConfig.PartitionIds = extractRuntimeInfo(hub, ctx, !d.CurrentConfig.SkipGetRuntimeInfo)
	return ctx, hub
}
