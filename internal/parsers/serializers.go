package parsers

import (
	"bytes"
	"encoding/gob"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
)

// TODO: Add summary
func SerializeOutboundMessage(m *d.OutboundMessage) []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(m)
	h.HandleError("Failed to Serialize OutboundMessage.", err, true)

	return res.Bytes()
}

// TODO: add summary
func DeserializeOutboundMessage(data []byte) *d.OutboundMessage {
	var m d.OutboundMessage
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&m)

	h.HandleError("Failed to deserialize OutboundMessage.", err, true)

	return &m
}

// TODO: Add summary
func SerializeInboundMessage(m *d.InboundMessage) []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(m)
	h.HandleError("Failed to Serialize InboundMessage.", err, true)

	return res.Bytes()
}

// TODO: add summary
func DeserializeInboundMessage(data []byte) *d.InboundMessage {
	var m d.InboundMessage
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&m)

	h.HandleError("Failed to deserialize InboundMessage.", err, true)

	return &m
}
