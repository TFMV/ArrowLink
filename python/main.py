import grpc
import pyarrow as pa
import pyarrow.ipc as ipc
from proto.dataexchange_pb2_grpc import ArrowDataServiceStub
from proto.dataexchange_pb2 import Empty

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = ArrowDataServiceStub(channel)
    
    response_stream = stub.GetArrowData(Empty())

    for response in response_stream:
        reader = ipc.RecordBatchStreamReader(pa.BufferReader(response.payload))
        table = reader.read_all()
        print(table)

if __name__ == '__main__':
    run()
