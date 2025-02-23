# ArrowLink

This is a simple implementation of a gRPC service that allows for the transfer of Arrow data between a Go server and a Python client.

## Overview

![Overview](docs/arrowlink.png)

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

## Notes

- The server and client are currently hardcoded to run on the same machine.
- The server and client are currently hardcoded to run on port 50051.

## Sample Output

```bash
Received DataFrame:
   id  value
0   1   3.14
```
