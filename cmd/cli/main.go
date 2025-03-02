package main

import (
	"fmt"
	"os"

	"github.com/TFMV/ArrowLink/arrow"
	"github.com/TFMV/ArrowLink/grpcserver"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "arrowlink",
	Short: "ArrowLink - High-performance data exchange between Go and Python",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the ArrowLink gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		rows, _ := cmd.Flags().GetInt("rows")

		logger, _ := zap.NewProduction()
		defer logger.Sync()

		arrowService := arrow.NewDemoArrowService(rows)
		grpcserver.StartGRPCServer(":"+port, logger, arrowService)
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Arrow data and output as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		rows, _ := cmd.Flags().GetInt("rows")
		output, _ := cmd.Flags().GetString("output")

		arrowService := arrow.NewDemoArrowService(rows)
		data, err := arrowService.GetData()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating data: %v\n", err)
			os.Exit(1)
		}

		// Convert Arrow data to JSON
		reader := arrow.NewArrowReader(data)
		df, err := reader.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting to JSON: %v\n", err)
			os.Exit(1)
		}

		if output == "" {
			// Print to stdout
			fmt.Println(string(df))
		} else {
			// Write to file
			if err := os.WriteFile(output, df, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Data written to %s\n", output)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(generateCmd)

	serverCmd.Flags().StringP("port", "p", "50051", "Port to listen on")
	serverCmd.Flags().IntP("rows", "r", 1000, "Number of rows to generate")

	generateCmd.Flags().IntP("rows", "r", 100, "Number of rows to generate")
	generateCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
