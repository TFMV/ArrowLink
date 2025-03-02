package arrow

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

// ArrowReader provides utilities for reading Arrow data
type ArrowReader struct {
	data []byte
	mem  memory.Allocator
}

// NewArrowReader creates a new ArrowReader from serialized Arrow data
func NewArrowReader(data []byte) *ArrowReader {
	return &ArrowReader{
		data: data,
		mem:  memory.NewGoAllocator(),
	}
}

// ToJSON converts Arrow data to JSON
func (r *ArrowReader) ToJSON() ([]byte, error) {
	reader, err := ipc.NewReader(bytes.NewReader(r.data), ipc.WithAllocator(r.mem))
	if err != nil {
		return nil, err
	}
	defer reader.Release()

	record, err := reader.Read()
	if err != nil {
		return nil, err
	}
	defer record.Release()

	// Convert to Go-friendly structure
	rows := make([]map[string]interface{}, record.NumRows())
	schema := record.Schema()

	for i := 0; i < int(record.NumRows()); i++ {
		row := make(map[string]interface{})
		for j, col := range record.Columns() {
			field := schema.Field(j)
			row[field.Name] = getValueAt(col, i)
		}
		rows[i] = row
	}

	return json.MarshalIndent(rows, "", "  ")
}

// getValueAt extracts a value from an Arrow array at the given index
func getValueAt(arr arrow.Array, i int) interface{} {
	if arr.IsNull(i) {
		return nil
	}

	switch arr := arr.(type) {
	case *array.Int64:
		return arr.Value(i)
	case *array.Float64:
		return arr.Value(i)
	case *array.String:
		return arr.Value(i)
	case *array.Boolean:
		return arr.Value(i)
	case *array.Timestamp:
		tsType := arr.DataType().(*arrow.TimestampType)
		return arr.Value(i).ToTime(tsType.Unit).Format(time.RFC3339)
	default:
		return "unsupported type"
	}
}
