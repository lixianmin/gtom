package mop

import (
	"log"
	"time"
)

type Options struct {
	After               TimestampGenerator // if nil defaults to gtom.LastOpTimestamp; not yet supported for ChangeStreamNS
	Filter              OpFilter           // op filter function that has access to type/ns/data
	NamespaceFilter     OpFilter           // op filter function that has access to type/ns ONLY
	OpLogDisabled       bool               // true to disable tailing the MongoDB oplog
	OpLogDatabaseName   string             // defaults to "local"
	OpLogCollectionName string             // defaults to "oplog.rs"
	ChannelSize         int                // defaults to 20
	BufferSize          int                // defaults to 50. used to batch fetch documents on bursts of activity
	BufferDuration      time.Duration      // defaults to 750 ms. after this timeout the batch is force fetched
	Ordering            OrderingGuarantee  // defaults to gtom.Oplog. ordering guarantee of events on the output channel as compared to the oplog
	WorkerCount         int                // defaults to 1. number of go routines batch fetching concurrently
	MaxWaitSecs         int                //
	UpdateDataAsDelta   bool               // set to true to only receive delta information in the Data field on updates (info straight from oplog)
	ChangeStreamNs      []string           // []string{"db.col1", "db.col2"}, MongoDB 3.6+ only; set to a slice to namespaces to read via MongoDB change streams
	DirectReadNs        []string           // []string{"db.users"}, set to a slice of namespaces (collections or views) to read data directly from
	DirectReadFilter    OpFilter           //
	DirectReadSplitMax  int32              // the max number of times to split a collection for concurrent reads (impacts memory consumption)
	DirectReadConcur    int                //
	DirectReadNoTimeout bool               //
	Unmarshal           DataUnmarshaller   //
	Pipe                PipelineBuilder    // an optional function to build aggregation pipelines
	PipeAllowDisk       bool               // true to allow MongoDB to use disk for aggregation pipeline options with large result sets
	Log                 *log.Logger        // pass your own logger
}

