package arrow

import (
	"bytes"

	arrow "github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

type ArrowService interface {
	GetData() ([]byte, error)
}

type arrowService struct {
	mem memory.Allocator
}

func NewArrowService() ArrowService {
	return &arrowService{
		mem: memory.NewGoAllocator(),
	}
}

func (s *arrowService) GetData() ([]byte, error) {
	// Define Arrow schema
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
	}, nil)

	// Create record builder
	builder := array.NewRecordBuilder(s.mem, schema)
	defer builder.Release()

	// Populate data
	intBuilder := builder.Field(0).(*array.Int64Builder)
	floatBuilder := builder.Field(1).(*array.Float64Builder)

	intBuilder.Append(1)
	floatBuilder.Append(3.14)

	record := builder.NewRecord()
	defer record.Release()

	// Serialize record using Arrow IPC
	var buf bytes.Buffer
	writer := ipc.NewWriter(&buf, ipc.WithSchema(schema))
	if err := writer.Write(record); err != nil {
		return nil, err
	}
	writer.Close()

	return buf.Bytes(), nil
}
