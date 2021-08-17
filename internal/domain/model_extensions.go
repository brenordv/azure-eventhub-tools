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
//  None.
//
// Receiver:
//  Instance of Message.
//
// Returns:
//  string representation of a Message.
func (m *InboundMessage) ToString() string {
	str := fmt.Sprintf(`---| DETAILS      |----------------------------------------------------------
id: %s
added to queue at: %s
event sequence number: %s
event offset: %s
Message processed at: %s
Filename: %s

---| MESSAGE BODY |----------------------------------------------------------
%s
---|          EOF |----------------------------------------------------------
`,
		m.EventId,
		m.QueuedTime.Format(time.RFC3339Nano),
		strconv.FormatInt(*m.EventSeqNumber, 10),
		strconv.FormatInt(*m.EventOffset, 10),
		m.ProcessedAt.Format(time.RFC3339Nano),
		m.SuggestedFilename,
		m.MsgData)

	return str
}

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