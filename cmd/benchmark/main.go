package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/TFMV/ArrowLink/arrow"
	"go.uber.org/zap"
)

func main() {
	// Parse command line flags
	minSize := flag.Int("min", 1000, "Minimum number of rows")
	maxSize := flag.Int("max", 1000000, "Maximum number of rows")
	steps := flag.Int("steps", 5, "Number of steps between min and max")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	fmt.Println("ArrowLink Benchmark")
	fmt.Println("==================")
	fmt.Printf("Testing data sizes from %d to %d rows in %d steps\n\n", *minSize, *maxSize, *steps)
	fmt.Println("Rows\tSize (KB)\tGen Time (ms)\tSer Time (ms)\tTotal Time (ms)")
	fmt.Println("----\t---------\t------------\t------------\t--------------")

	// Calculate step size
	stepSize := (*maxSize - *minSize) / (*steps - 1)
	if stepSize <= 0 {
		stepSize = 1
	}

	for size := *minSize; size <= *maxSize; size += stepSize {
		// Create service with specific size
		service := arrow.NewDemoArrowService(size)

		// Measure time
		start := time.Now()
		data, err := service.GetData()
		if err != nil {
			log.Fatalf("Error generating data: %v", err)
		}
		elapsed := time.Since(start)

		// Get metrics from the benchmark service
		metrics := service.(*arrow.DemoArrowService).GetMetrics()

		fmt.Printf("%d\t%d\t\t%.2f\t\t%.2f\t\t%.2f\n",
			size,
			len(data)/1024,
			float64(metrics.GenerationTime)/float64(time.Millisecond),
			float64(metrics.SerializationTime)/float64(time.Millisecond),
			float64(elapsed)/float64(time.Millisecond))
	}
}
