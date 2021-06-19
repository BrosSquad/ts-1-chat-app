import 'package:chat_app/proto/chat.pbgrpc.dart';
import 'package:fixnum/fixnum.dart';
import 'package:grpc/grpc.dart';

class ChatService {
  late ChatClient _client;

  ChatService(ClientChannel channel) {
    _client = ChatClient(channel);
  }

  Stream<MessageResponse> connect(Int64 userID) async* {
    try {
      ResponseStream<MessageResponse> messagesStream =
          _client.connect(ConnectRequest(userId: userID));
      await for (var message in messagesStream) {
        // print('Received message $message');
        yield message;
      }
    } catch (error) {
      throw 'Could not connect to message stream $error';
    }
  }

  Future<void> sendMessage(
      {required Int64 userID, required String text}) async {
    try {
      await _client.sendMessage(
        MessageRequest(
          userId: userID,
          text: text,
        ),
      );
      print('Message sent $text');
    } catch (error) {
      print('Error sending message $error');
    }
  }
}
