import grpc
import pyarrow as pa
import pyarrow.ipc as ipc
import logging
import time
import argparse
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from grpc import RpcError

from proto.dataexchange_pb2_grpc import ArrowDataServiceStub
from proto.dataexchange_pb2 import Empty


class LoggingInterceptor(
    grpc.UnaryUnaryClientInterceptor,
    grpc.UnaryStreamClientInterceptor,
    grpc.StreamUnaryClientInterceptor,
    grpc.StreamStreamClientInterceptor,
):
    def intercept_unary_unary(self, continuation, client_call_details, request):
        logging.info(f"Unary call to {client_call_details.method}")
        return continuation(client_call_details, request)

    def intercept_unary_stream(self, continuation, client_call_details, request):
        logging.info(f"Unary-stream call to {client_call_details.method}")
        return continuation(client_call_details, request)

    def intercept_stream_unary(
        self, continuation, client_call_details, request_iterator
    ):
        logging.info(f"Stream-unary call to {client_call_details.method}")
        return continuation(client_call_details, request_iterator)

    def intercept_stream_stream(
        self, continuation, client_call_details, request_iterator
    ):
        logging.info(f"Stream-stream call to {client_call_details.method}")
        return continuation(client_call_details, request_iterator)


def get_secure_channel(target, ca_cert_file, options):
    """
    Creates a secure gRPC channel using the provided CA certificate.
    """
    try:
        with open(ca_cert_file, "rb") as f:
            trusted_certs = f.read()
        credentials = grpc.ssl_channel_credentials(root_certificates=trusted_certs)
        return grpc.secure_channel(target, credentials, options=options)
    except FileNotFoundError:
        logging.warning(
            f"Certificate file {ca_cert_file} not found. Using insecure channel."
        )
        return grpc.insecure_channel(target, options=options)
    except Exception as e:
        logging.error(f"Failed to create secure channel: {e}")
        logging.warning(
            "Falling back to insecure channel (not recommended in production)."
        )
        return grpc.insecure_channel(target, options=options)


def run():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description="ArrowLink Python Client")
    parser.add_argument(
        "--benchmark", action="store_true", help="Run performance benchmark"
    )
    parser.add_argument(
        "--visualize", action="store_true", help="Generate visualization"
    )
    parser.add_argument(
        "--cert", type=str, default="certs/ca.crt", help="Path to CA certificate"
    )
    args = parser.parse_args()

    logging.basicConfig(level=logging.INFO)
    target = "localhost:50051"
    ca_cert_file = args.cert

    # Channel options to tune message sizes and keep-alive for high scale.
    options = [
        ("grpc.max_send_message_length", 50 * 1024 * 1024),
        ("grpc.max_receive_message_length", 50 * 1024 * 1024),
        ("grpc.keepalive_time_ms", 10000),
        ("grpc.keepalive_timeout_ms", 5000),
    ]

    # Create channel (secure if possible, otherwise insecure)
    channel = get_secure_channel(target, ca_cert_file, options)

    # Wrap the channel with a logging interceptor.
    intercepted_channel = grpc.intercept_channel(channel, LoggingInterceptor())
    stub = ArrowDataServiceStub(intercepted_channel)

    max_retries = 3
    retry_delay = 5  # seconds

    # Benchmark mode
    if args.benchmark:
        start_time = time.time()

    # Implement retry logic with a call deadline.
    for attempt in range(1, max_retries + 1):
        try:
            logging.info("Calling GetArrowData (attempt %d)...", attempt)
            # Set a deadline of 30 seconds for the RPC call.
            response_stream = stub.GetArrowData(Empty(), timeout=30)
            for response in response_stream:
                try:
                    reader = ipc.RecordBatchStreamReader(
                        pa.BufferReader(response.payload)
                    )
                    table = reader.read_all()
                    df = table.to_pandas()

                    if args.benchmark:
                        end_time = time.time()
                        logging.info(
                            f"Received {len(df)} rows in {end_time - start_time:.4f} seconds"
                        )
                        logging.info(
                            f"Throughput: {len(df) / (end_time - start_time):.2f} rows/second"
                        )
                    else:
                        logging.info(
                            f"Received Arrow table with {len(df)} rows and {len(df.columns)} columns"
                        )
                        logging.info(f"Schema: {table.schema}")
                        logging.info(f"Sample data:\n{df.head()}")

                    if args.visualize:
                        visualize_data(df)

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


def visualize_data(df):
    """Generate visualizations for the received data"""
    # Set the style
    sns.set(style="whitegrid")

    # Create a figure with multiple subplots
    fig, axes = plt.subplots(2, 2, figsize=(15, 10))

    # Plot 1: Time series of values
    df.plot(x="timestamp", y="value", ax=axes[0, 0], title="Time Series of Values")

    # Plot 2: Distribution of values by category
    sns.boxplot(x="category", y="value", data=df, ax=axes[0, 1])
    axes[0, 1].set_title("Value Distribution by Category")

    # Plot 3: Count by category
    category_counts = df["category"].value_counts()
    category_counts.plot.bar(ax=axes[1, 0], title="Count by Category")

    # Plot 4: Valid vs Invalid counts
    valid_counts = df["is_valid"].value_counts()
    valid_counts.plot.pie(
        ax=axes[1, 1], autopct="%1.1f%%", title="Valid vs Invalid Records"
    )

    plt.tight_layout()
    plt.savefig("arrowlink_visualization.png")
    logging.info("Visualization saved to arrowlink_visualization.png")
    plt.show()


if __name__ == "__main__":
    run()
