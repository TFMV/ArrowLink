# ArrowLink

ArrowLink is a demonstration of using Apache Arrow for efficient data exchange between a Go-based gRPC server and a Python client. It leverages Arrow’s zero-copy serialization and gRPC streaming to transfer structured data at high speeds.

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

    %% Modern Color Palette
    style A fill:#E6F3FF,stroke:#1E90FF,stroke-width:2px,color:#1E90FF %% Light blue, vibrant stroke
    style B fill:#F0F8FF,stroke:#4682B4,stroke-width:1px %% Softer blue
    style C fill:#F0F8FF,stroke:#4682B4,stroke-width:1px
    style L1 fill:#B3D9FF,stroke:#1E90FF,stroke-width:2px %% Medium blue
    style SC fill:#CCE5FF,stroke:#1E90FF,stroke-width:2px %% Light-medium blue
    style D fill:#FFD700,stroke:#DAA520,stroke-width:2px,color:#DAA520 %% Gold for gRPC Channel

    style E fill:#F5E6FF,stroke:#8A2BE2,stroke-width:2px,color:#8A2BE2 %% Light purple for proto
    style F fill:#E6CCFF,stroke:#8A2BE2,stroke-width:1px

    style G fill:#E6FFE6,stroke:#32CD32,stroke-width:2px,color:#32CD32 %% Light green for Go entry
    style M fill:#CCFFCC,stroke:#32CD32,stroke-width:1px %% Soft green
    style S fill:#B3FFB3,stroke:#228B22,stroke-width:2px %% Medium green
    style GS fill:#99FF99,stroke:#228B22,stroke-width:2px %% Light green

    style MW fill:#FFFFE6,stroke:#FFD700,stroke-width:1px %% Pale yellow
    style REC fill:#FFFACD,stroke:#FFD700,stroke-width:1px
    style H fill:#FFFFCC,stroke:#FFD700,stroke-width:2px
    style N fill:#FFFFCC,stroke:#FFD700,stroke-width:1px

    style L2 fill:#FFE6E6,stroke:#FF4500,stroke-width:2px,color:#FF4500 %% Light red for Arrow
    style O fill:#FFD9D9,stroke:#FF4500,stroke-width:1px
    style P fill:#FFD9D9,stroke:#FF4500,stroke-width:1px
    style I fill:#FFCCCC,stroke:#FF4500,stroke-width:2px
    style Q fill:#FFB3B3,stroke:#FF4500,stroke-width:2px

    style J fill:#FFF0E6,stroke:#FF6347,stroke-width:2px %% Soft coral
    style R fill:#FFE6CC,stroke:#FF6347,stroke-width:2px
```

## Key Features

- 🚀 gRPC-based communication between Go and Python
- 🔄 Apache Arrow for efficient binary data exchange
- 📡 Streaming support for handling large datasets
- 🏎 High-speed, zero-copy serialization for optimal performance
- 🔧 Extensible architecture for integrating with real-world data systems

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

## Sample Output

```bash
pyarrow.Table
id: int64 not null
value: double not null
----
id: [[1]]
value: [[3.14]]
```
