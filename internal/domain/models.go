package domain

import "time"

// Config is the configuration read from the file passed via command line argument.
type Config struct {
	EventhubConnectionString string         `json:"eventhubConnString"`
	EntityPath               string         `json:"entityPath"`
	SkipGetRuntimeInfo       bool           `json:"skipGetRuntimeInfo"`
	BadgerConfig             BadgerConfig   `json:"badgerConfig"`
	InboundConfig            InboundConfig  `json:"inboundConfig"`
	OutboundConfig           OutboundConfig `json:"outboundConfig"`
	PartitionIds             []string
}

type OutboundConfig struct {
	OutboundFolder   string `json:"outboundFolder"`
	Buffered         bool   `json:"buffered"`
	JustSendBuffered bool   `json:"justSendBuffered"`
	IgnoreStatus     bool   `json:"ignoreStatus"`
}

type InboundConfig struct {
	ConsumerGroup    string   `json:"consumerGroup"`
	PartitionId      int      `json:"partitionId"`
	InboundFolder    string   `json:"inboundFolder"`
	ReadToFile       bool     `json:"readToFile"`
	IgnoreCheckpoint bool     `json:"ignoreCheckpoint"`
	DumpContentOnly  bool     `json:"dumpContentOnly"`
	DumpFilter       []string `json:"dumpFilter"`
}

// TODO: add summary
type BadgerConfig struct {
	Verbose                    bool   `json:"verboseMode"`
	BadgerSkipCompactL0OnClose bool   `json:"badgerSkipCompactL0OnClose"`
	BadgerValueLogFileSize     int64  `json:"badgerValueLogFileSize"`
	OutboundBaseDir            string `json:"outboundBaseDir"`
	OutboundDir                string `json:"outboundDir"`
	OutboundValueDir           string `json:"outboundValueDir"`
	InboundBaseDir             string `json:"inboundBaseDir"`
	InboundDir                 string `json:"inboundDir"`
	InboundValueDir            string `json:"inboundValueDir"`
}

// TODO: add summary
type FuncReturnError func() error

// TODO: add summary
type FuncReturnErrorIntArg func(n int) error

// TODO: add summary
type OutboundMsgStatus int

const (
	New OutboundMsgStatus = iota
	Buffered
	Sent
	Error
)

type InboundMsgStatus int

const (
	Read InboundMsgStatus = iota
	Exported
)

// TODO: add summary
type OutboundMessage struct {
	Id           []byte
	FullFilename string
	Content      string
	ProcessedAt  time.Time
	Status       OutboundMsgStatus
}

// TODO: add summary
type InboundMessage struct {
	Id                []byte
	EventId           string
	PartitionKey      *string
	PartitionId       string
	QueuedTime        time.Time
	EventSeqNumber    *int64
	EventOffset       *int64
	SuggestedFilename string
	ProcessedAt       time.Time
	MsgData           string
	Status            InboundMsgStatus
}
