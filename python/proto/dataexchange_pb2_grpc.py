# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

from . import dataexchange_pb2 as dataexchange__pb2

GRPC_GENERATED_VERSION = '1.70.0'
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in dataexchange_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
    )


class ArrowDataServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetArrowData = channel.unary_stream(
                '/dataexchange.ArrowDataService/GetArrowData',
                request_serializer=dataexchange__pb2.Empty.SerializeToString,
                response_deserializer=dataexchange__pb2.ArrowData.FromString,
                _registered_method=True)
        self.SendArrowData = channel.stream_unary(
                '/dataexchange.ArrowDataService/SendArrowData',
                request_serializer=dataexchange__pb2.ArrowData.SerializeToString,
                response_deserializer=dataexchange__pb2.Ack.FromString,
                _registered_method=True)


class ArrowDataServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def GetArrowData(self, request, context):
        """Streaming response for efficient data transfer
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SendArrowData(self, request_iterator, context):
        """Accepts Arrow data and processes it
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ArrowDataServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetArrowData': grpc.unary_stream_rpc_method_handler(
                    servicer.GetArrowData,
                    request_deserializer=dataexchange__pb2.Empty.FromString,
                    response_serializer=dataexchange__pb2.ArrowData.SerializeToString,
            ),
            'SendArrowData': grpc.stream_unary_rpc_method_handler(
                    servicer.SendArrowData,
                    request_deserializer=dataexchange__pb2.ArrowData.FromString,
                    response_serializer=dataexchange__pb2.Ack.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'dataexchange.ArrowDataService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers('dataexchange.ArrowDataService', rpc_method_handlers)


 # This class is part of an EXPERIMENTAL API.
class ArrowDataService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def GetArrowData(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(
            request,
            target,
            '/dataexchange.ArrowDataService/GetArrowData',
            dataexchange__pb2.Empty.SerializeToString,
            dataexchange__pb2.ArrowData.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def SendArrowData(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(
            request_iterator,
            target,
            '/dataexchange.ArrowDataService/SendArrowData',
            dataexchange__pb2.ArrowData.SerializeToString,
            dataexchange__pb2.Ack.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
