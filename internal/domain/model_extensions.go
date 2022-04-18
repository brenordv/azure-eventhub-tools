package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ToString converts a Message to string.
//
// Parameters:
//  None. Uses receiver data.
//
// Receiver:
//  Instance of Message.
//
// Returns:
//  string representation of a Message.
func (m *InboundMessage) ToString() string {
	var partitionKey string
	if m.PartitionKey != nil {
		partitionKey = *m.PartitionKey
	} else {
		partitionKey = ""
	}

	str := fmt.Sprintf(`---| DETAILS      |----------------------------------------------------------
id: %s
partition key: %s
added to queue at: %s
partition id: %s
event sequence number: %s
event offset: %s
Message processed at: %s
Filename: %s

---| MESSAGE BODY |----------------------------------------------------------
%s
---|          EOF |----------------------------------------------------------
`,
		m.EventId,
		partitionKey,
		m.QueuedTime.Format(time.RFC3339Nano),
		m.PartitionId,
		strconv.FormatInt(*m.EventSeqNumber, 10),
		strconv.FormatInt(*m.EventOffset, 10),
		m.ProcessedAt.Format(time.RFC3339Nano),
		m.SuggestedFilename,
		m.MsgData)

	return str
}

// ContentHasFilterKeywords checks if the string c contains any of the filters used.
//
// Parameters:
//  c: text to be checked.
//
// Receiver:
//  Instance of InboundConfig
//
// Returns:
//  true if c contains any of the filters informed or if there are no filters to use. false otherwise.
func (ic *InboundConfig) ContentHasFilterKeywords(c string) bool {
	if CurrentConfig.InboundConfig.DumpFilter == nil {
		return false
	}
	c = strings.ToLower(c)
	for _, f := range CurrentConfig.InboundConfig.DumpFilter {
		if strings.Contains(c, f) {
			return true
		}
	}
	return false
}
