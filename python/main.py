import grpc
import pyarrow as pa
import pyarrow.ipc as ipc
import logging
import time
from grpc import RpcError

from proto.dataexchange_pb2_grpc import ArrowDataServiceStub
from proto.dataexchange_pb2 import Empty

class LoggingInterceptor(
    grpc.UnaryUnaryClientInterceptor,
    grpc.UnaryStreamClientInterceptor,
    grpc.StreamUnaryClientInterceptor,
    grpc.StreamStreamClientInterceptor
):
    def intercept_unary_unary(self, continuation, client_call_details, request):
        logging.info(f"Unary call to {client_call_details.method}")
        return continuation(client_call_details, request)

    def intercept_unary_stream(self, continuation, client_call_details, request):
        logging.info(f"Unary-stream call to {client_call_details.method}")
        return continuation(client_call_details, request)

    def intercept_stream_unary(self, continuation, client_call_details, request_iterator):
        logging.info(f"Stream-unary call to {client_call_details.method}")
        return continuation(client_call_details, request_iterator)

    def intercept_stream_stream(self, continuation, client_call_details, request_iterator):
        logging.info(f"Stream-stream call to {client_call_details.method}")
        return continuation(client_call_details, request_iterator)

def get_secure_channel(target, ca_cert_file, options):
    """
    Creates a secure gRPC channel using the provided CA certificate.
    """
    with open(ca_cert_file, 'rb') as f:
        trusted_certs = f.read()
    credentials = grpc.ssl_channel_credentials(root_certificates=trusted_certs)
    return grpc.secure_channel(target, credentials, options=options)

def run():
    logging.basicConfig(level=logging.INFO)
    target = 'localhost:50051'
    ca_cert_file = 'ArrowLink/certs/ca.crt'

    # Channel options to tune message sizes and keep-alive for high scale.
    options = [
        ('grpc.max_send_message_length', 50 * 1024 * 1024),
        ('grpc.max_receive_message_length', 50 * 1024 * 1024),
        ('grpc.keepalive_time_ms', 10000),
        ('grpc.keepalive_timeout_ms', 5000),
    ]
    
    # Attempt to create a secure channel.
    try:
        channel = get_secure_channel(target, ca_cert_file, options)
        logging.info("Secure channel established.")
    except Exception as e:
        logging.error("Failed to create secure channel. Falling back to insecure channel (not recommended in production).", exc_info=e)
        channel = grpc.insecure_channel(target, options=options)
    
    # Wrap the channel with a logging interceptor.
    intercepted_channel = grpc.intercept_channel(channel, LoggingInterceptor())
    stub = ArrowDataServiceStub(intercepted_channel)
    
    max_retries = 3
    retry_delay = 5  # seconds
    
    # Implement retry logic with a call deadline.
    for attempt in range(1, max_retries + 1):
        try:
            logging.info("Calling GetArrowData (attempt %d)...", attempt)
            # Set a deadline of 30 seconds for the RPC call.
            response_stream = stub.GetArrowData(Empty(), timeout=30)
            for response in response_stream:
                try:
                    reader = ipc.RecordBatchStreamReader(pa.BufferReader(response.payload))
                    table = reader.read_all()
                    logging.info("Received Arrow table:\n%s", table)
                except Exception as data_err:
                    logging.error("Error processing Arrow data.", exc_info=data_err)
            break
        except RpcError as rpc_err:
            logging.error("gRPC error on attempt %d: %s", attempt, rpc_err)
            if attempt < max_retries:
                logging.info("Retrying in %d seconds...", retry_delay)
                time.sleep(retry_delay)
            else:
                logging.error("Max retries exceeded. Exiting.")
    channel.close()

if __name__ == '__main__':
    run()
