package arrow

import (
	"bytes"
	"math/rand"
	"time"

	arrow "github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

type Metrics struct {
	GenerationTime    time.Duration
	SerializationTime time.Duration
}

type DemoArrowService struct {
	mem      memory.Allocator
	dataSize int
	metrics  Metrics
}

func NewDemoArrowService(dataSize int) ArrowService {
	return &DemoArrowService{
		mem:      memory.NewGoAllocator(),
		dataSize: dataSize,
	}
}

func (s *DemoArrowService) GetMetrics() Metrics {
	return s.metrics
}

func (s *DemoArrowService) GetData() ([]byte, error) {
	startGen := time.Now()

	// Define a more complex Arrow schema with various data types
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "timestamp", Type: arrow.FixedWidthTypes.Timestamp_ms},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		{Name: "category", Type: arrow.BinaryTypes.String},
		{Name: "is_valid", Type: arrow.FixedWidthTypes.Boolean},
	}, nil)

	// Create record builder
	builder := array.NewRecordBuilder(s.mem, schema)
	defer builder.Release()

	// Get builders for each field
	idBuilder := builder.Field(0).(*array.Int64Builder)
	tsBuilder := builder.Field(1).(*array.TimestampBuilder)
	valueBuilder := builder.Field(2).(*array.Float64Builder)
	categoryBuilder := builder.Field(3).(*array.StringBuilder)
	validBuilder := builder.Field(4).(*array.BooleanBuilder)

	// Generate random data
	rand.Seed(time.Now().UnixNano())
	categories := []string{"A", "B", "C", "D", "E"}
	now := time.Now()

	// Populate with random data
	for i := 0; i < s.dataSize; i++ {
		idBuilder.Append(int64(i))
		tsBuilder.Append(arrow.Timestamp(now.Add(time.Duration(i)*time.Second).UnixNano() / int64(time.Millisecond)))
		valueBuilder.Append(rand.Float64() * 100)
		categoryBuilder.Append(categories[rand.Intn(len(categories))])
		validBuilder.Append(rand.Intn(10) > 2) // 70% valid
	}

	record := builder.NewRecord()
	defer record.Release()

	// After record creation
	s.metrics.GenerationTime = time.Since(startGen)

	// Before serialization
	startSer := time.Now()

	// Serialize record using Arrow IPC
	var buf bytes.Buffer
	writer := ipc.NewWriter(&buf, ipc.WithSchema(schema))
	if err := writer.Write(record); err != nil {
		return nil, err
	}
	writer.Close()

	// After serialization
	s.metrics.SerializationTime = time.Since(startSer)

	return buf.Bytes(), nil
}
