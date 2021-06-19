import 'package:chat_app/services/auth.dart';
import 'package:chat_app/services/chat.dart';
import 'package:flutter/widgets.dart';
import 'package:grpc/grpc.dart';

ClientChannel channel = ClientChannel(
  "localhost",
  port: 3000,
  options: const ChannelOptions(
    credentials: ChannelCredentials.insecure(),
  ),
);

class GRPCClient extends InheritedWidget {
  final AuthService authService = AuthService(channel);
  final ChatService chatService = ChatService(channel);

  GRPCClient({required Widget child}) : super(child: child);

  @override
  bool updateShouldNotify(covariant InheritedWidget oldWidget) {
    return true;
  }

  static GRPCClient of(BuildContext context) {
    final GRPCClient? result =
        context.dependOnInheritedWidgetOfExactType<GRPCClient>();
    assert(result != null, 'No GRPCClient found in context');
    return result!;
  }
}
