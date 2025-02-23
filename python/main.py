import grpc
import pyarrow as pa
import pyarrow.ipc as ipc
import pandas as pd
from proto.dataexchange_pb2_grpc import ArrowDataServiceStub
from proto.dataexchange_pb2 import ArrowData, Empty  # Add any other message types you need

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = ArrowDataServiceStub(channel)
    
    response_stream = stub.GetArrowData(Empty())
    
    for response in response_stream:
        # Deserialize the Arrow IPC stream
        reader = ipc.RecordBatchStreamReader(pa.BufferReader(response.payload))
        table = reader.read_all()
        # Convert to Pandas DataFrame for further analysis
        df = table.to_pandas()
        print("Received DataFrame:")
        print(df)

if __name__ == '__main__':
    run()
