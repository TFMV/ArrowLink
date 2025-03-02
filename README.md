# ArrowLink

ArrowLink is a demonstration of using Apache Arrow for efficient data exchange between a Go-based gRPC server and a Python client. It leverages Arrow's zero-copy serialization and gRPC streaming to transfer structured data at high speeds.

## Overview

```mermaid
graph TD
    subgraph "Python Client"
        A[main.py] --> |imports| B[proto/dataexchange_pb2_grpc.py]
        A --> |imports| C[proto/dataexchange_pb2.py]
        B --> |relative import| C
        A --> |creates| L1[LoggingInterceptor]
        A --> |creates| SC[Secure Channel]
        SC --> |wraps| L1
        L1 --> |"stub.GetArrowData(Empty())"| D[gRPC Channel]
    end

    subgraph "Proto Definition"
        E[dataexchange.proto] --> |protoc| B
        E --> |protoc| C
        E --> |protoc| F[Go: dataexchange.pb.go]
    end

    subgraph "Go Server Components"
        G[arrowlink.go] --> |initializes| M[Zap Logger]
        G --> |creates| S[Arrow Service]
        G --> |starts| GS[gRPC Server]
    end

    subgraph "gRPC Server"
        GS --> |uses| MW[Middleware]
        MW --> |logging| M
        MW --> |recovery| REC[Recovery Handler]
        GS --> |implements| H[ArrowDataService]
        H --> |embedded| N[ArrowDataServiceServer]
        H --> |uses| S
    end

    subgraph "Arrow Service"
        S --> |imports| L2[arrow-go/v18]
        S --> |"Create Schema"| O[Arrow Schema]
        O --> |"Build Record"| P[Record Builder]
        P --> |"Append Data"| I[Arrow Record Batch]
        I --> |"IPC Writer"| Q[Buffer]
    end

    D <--> |"gRPC (TLS/50051)"| GS
    Q --> |"Serialized bytes"| D
    D --> |"response.payload"| J[PyArrow Reader]
    J --> |"read_all()"| R[PyArrow Table]

    %% Python Client section - Blue theme
    style A fill:#2B6CB0,stroke:#1A365D,stroke-width:2px,color:#FFFFFF
    style B fill:#4299E1,stroke:#2C5282,stroke-width:1px,color:#FFFFFF
    style C fill:#4299E1,stroke:#2C5282,stroke-width:1px,color:#FFFFFF
    style L1 fill:#3182CE,stroke:#2A4365,stroke-width:2px,color:#FFFFFF
    style SC fill:#2B6CB0,stroke:#1A365D,stroke-width:2px,color:#FFFFFF

    %% Proto Definition - Purple theme
    style E fill:#6B46C1,stroke:#44337A,stroke-width:2px,color:#FFFFFF
    style F fill:#805AD5,stroke:#553C9A,stroke-width:1px,color:#FFFFFF

    %% Go Server - Green theme
    style G fill:#2F855A,stroke:#22543D,stroke-width:2px,color:#FFFFFF
    style M fill:#38A169,stroke:#276749,stroke-width:1px,color:#FFFFFF
    style S fill:#48BB78,stroke:#2F855A,stroke-width:2px,color:#FFFFFF
    style GS fill:#68D391,stroke:#38A169,stroke-width:2px,color:#000000

    %% gRPC Server - Orange theme
    style MW fill:#DD6B20,stroke:#9C4221,stroke-width:1px,color:#FFFFFF
    style REC fill:#ED8936,stroke:#C05621,stroke-width:1px,color:#FFFFFF
    style H fill:#F6AD55,stroke:#DD6B20,stroke-width:2px,color:#000000
    style N fill:#FBD38D,stroke:#ED8936,stroke-width:1px,color:#000000

    %% Arrow Service - Red theme
    style L2 fill:#C53030,stroke:#822727,stroke-width:2px,color:#FFFFFF
    style O fill:#E53E3E,stroke:#C53030,stroke-width:1px,color:#FFFFFF
    style P fill:#E53E3E,stroke:#C53030,stroke-width:1px,color:#FFFFFF
    style I fill:#F56565,stroke:#E53E3E,stroke-width:2px,color:#FFFFFF
    style Q fill:#FC8181,stroke:#F56565,stroke-width:2px,color:#000000

    %% Communication - Yellow theme
    style D fill:#D69E2E,stroke:#975A16,stroke-width:2px,color:#000000
    style J fill:#ECC94B,stroke:#D69E2E,stroke-width:2px,color:#000000
    style R fill:#F6E05E,stroke:#ECC94B,stroke-width:2px,color:#000000
```

## Key Features

- 🚀 gRPC-based communication between Go and Python
- 🔄 Apache Arrow for efficient binary data exchange
- 📡 Streaming support for handling large datasets
- 🏎 High-speed, zero-copy serialization for optimal performance
- 🔧 Extensible architecture for integrating with real-world data systems

## Performance Benchmarks

ArrowLink demonstrates exceptional performance when handling large datasets. Below are the results from our benchmark tests:

| Rows      | Size (KB) | Generation Time (ms) | Serialization Time (ms) | Total Time (ms) |
| --------- | --------- | -------------------- | ----------------------- | --------------- |
| 1,000     | 29        | 0.77                 | 0.17                    | 0.95            |
| 250,750   | 7,132     | 27.83                | 3.29                    | 31.12           |
| 500,500   | 14,236    | 35.62                | 3.51                    | 39.14           |
| 750,250   | 21,339    | 45.42                | 6.50                    | 51.93           |
| 1,000,000 | 28,443    | 52.50                | 3.58                    | 56.09           |

These results highlight ArrowLink's ability to:

- Process 1 million rows in just 56ms
- Maintain efficient serialization even as data size increases
- Achieve compression ratios that keep data sizes manageable

## Use Cases

ArrowLink is useful in scenarios where high-performance, structured data exchange is required across multiple programming environments. Some practical applications include:

1. Real-Time Data Pipelines
   - Send structured data from Go-based ingestion services to Python-based analytics engines.
   - Example: Streaming sensor data from an IoT device to a Python ML model.
2. Machine Learning Inference
   - Use Go to handle API requests while forwarding Arrow-encoded data to a Python ML inference engine.
   - Example: A recommendation system where Go receives user queries and Python processes the embeddings.
3. ETL & Data Processing Workflows
   - Move large datasets between Go and Python without expensive JSON serialization/deserialization.
   - Example: A Go service collects logs and sends them to Python for batch processing with Pandas.
4. Vector Search & AI Pipelines
   - Utilize Arrow for fast embedding transmission between Go-based search engines and Python-based vector search libraries (e.g., FAISS, Annoy).
   - Example: A document similarity engine where Go serves the front-end API and Python handles vector indexing.

## Setup

### Install the dependencies

```bash
pip install -r python/requirements.txt
```

### Run the server

```bash
go run arrowlink.go
```

### Run the client

```bash
python python/main.py
```

### Run the benchmark

```bash
go run cmd/benchmark/main.go --min 1000 --max 1000000 --steps 5
```

### Run the dashboard

```bash
streamlit run python/dashboard.py
```

## Sample Output

```bash
pyarrow.Table
id: int64 not null
value: double not null
----
id: [[1]]
value: [[3.14]]
```

## Docker Support

You can also run ArrowLink using Docker Compose:

```bash
docker-compose up
```

This will start both the server and the dashboard, making them accessible at:

- gRPC Server: localhost:50051
- Dashboard: <http://localhost:8501>
