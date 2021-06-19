///
//  Generated code. Do not modify.
//  source: auth.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,unnecessary_const,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type,unnecessary_this,prefer_final_fields

import 'dart:async' as $async;

import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'auth.pb.dart' as $1;
export 'auth.pb.dart';

class AuthClient extends $grpc.Client {
  static final _$register =
      $grpc.ClientMethod<$1.RegisterRequest, $1.RegisterResponse>(
          '/auth.Auth/Register',
          ($1.RegisterRequest value) => value.writeToBuffer(),
          ($core.List<$core.int> value) =>
              $1.RegisterResponse.fromBuffer(value));

  AuthClient($grpc.ClientChannel channel,
      {$grpc.CallOptions? options,
      $core.Iterable<$grpc.ClientInterceptor>? interceptors})
      : super(channel, options: options, interceptors: interceptors);

  $grpc.ResponseFuture<$1.RegisterResponse> register($1.RegisterRequest request,
      {$grpc.CallOptions? options}) {
    return $createUnaryCall(_$register, request, options: options);
  }
}

abstract class AuthServiceBase extends $grpc.Service {
  $core.String get $name => 'auth.Auth';

  AuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$1.RegisterRequest, $1.RegisterResponse>(
        'Register',
        register_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.RegisterRequest.fromBuffer(value),
        ($1.RegisterResponse value) => value.writeToBuffer()));
  }

  $async.Future<$1.RegisterResponse> register_Pre(
      $grpc.ServiceCall call, $async.Future<$1.RegisterRequest> request) async {
    return register(call, await request);
  }

  $async.Future<$1.RegisterResponse> register(
      $grpc.ServiceCall call, $1.RegisterRequest request);
}
