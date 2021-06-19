///
//  Generated code. Do not modify.
//  source: chat.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,unnecessary_const,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type,unnecessary_this,prefer_final_fields

import 'dart:async' as $async;

import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'chat.pb.dart' as $0;
export 'chat.pb.dart';

class ChatClient extends $grpc.Client {
  static final _$sendMessage = $grpc.ClientMethod<$0.MessageRequest, $0.Empty>(
      '/chat.Chat/SendMessage',
      ($0.MessageRequest value) => value.writeToBuffer(),
      ($core.List<$core.int> value) => $0.Empty.fromBuffer(value));
  static final _$connect =
      $grpc.ClientMethod<$0.ConnectRequest, $0.MessageResponse>(
          '/chat.Chat/Connect',
          ($0.ConnectRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $0.MessageResponse.fromBuffer(value));

  ChatClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options, interceptors: interceptors);

  $grpc.ResponseFuture<$0.Empty> sendMessage($0.MessageRequest request,
      {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$sendMessage, request, options: options);
  }

  $grpc.ResponseStream<$0.MessageResponse> connect($0.ConnectRequest request,
      {$grpc.CallOptions? options}) {
    return $createStreamingCall(
        _$connect, $async.Stream.fromIterable([request]),
        options: options);
  }
}

abstract class ChatServiceBase extends $grpc.Service {
  $core.String get $name => 'chat.Chat';

  ChatServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.MessageRequest, $0.Empty>(
        'SendMessage',
        sendMessage_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.MessageRequest.fromBuffer(value),
        ($0.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConnectRequest, $0.MessageResponse>(
        'Connect',
        connect_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.ConnectRequest.fromBuffer(value),
        ($0.MessageResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.Empty> sendMessage_Pre(
      $grpc.ServiceCall call, $async.Future<$0.MessageRequest> request) async {
    return sendMessage(call, await request);
  }

  $async.Stream<$0.MessageResponse> connect_Pre(
      $grpc.ServiceCall call, $async.Future<$0.ConnectRequest> request) async* {
    yield* connect(call, await request);
  }

  $async.Future<$0.Empty> sendMessage(
      $grpc.ServiceCall call, $0.MessageRequest request);
  $async.Stream<$0.MessageResponse> connect(
      $grpc.ServiceCall call, $0.ConnectRequest request);
}
