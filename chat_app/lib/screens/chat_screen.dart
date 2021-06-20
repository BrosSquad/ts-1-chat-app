import 'package:chat_app/proto/chat.pb.dart';
import 'package:chat_app/services/chat.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:timeago/timeago.dart' as timeago;

import '../grpc_client.dart';

class ChatScreen extends StatefulWidget {
  const ChatScreen({Key? key, required this.service, required this.userID})
      : super(key: key);

  final ChatService service;
  final Int64 userID;

  @override
  _ChatScreenState createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  late TextEditingController _controller;
  List<MessageResponse> messages = [];
  late Stream<MessageResponse> stream;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController();
    this.stream = widget.service.connect(widget.userID);
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> sendMessage({required Int64 userID}) async {
    await widget.service
        .sendMessage(userID: userID, text: _controller.value.text);
    _controller.clear();
  }

  @override
  Widget build(BuildContext context) {
    String username = GRPCClient.of(context).authService.username;

    return Scaffold(
      appBar: AppBar(
        title: Text(username),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => sendMessage(userID: widget.userID),
        child: Icon(Icons.send),
      ),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: StreamBuilder<MessageResponse>(
                  stream: this.stream,
                  builder: (context, snapshot) {
                    print(snapshot.requireData);
                    if (!snapshot.hasData) {
                      return Center(
                        child: CircularProgressIndicator(),
                      );
                    }
                    messages.add(snapshot.data!);
                    print(messages.length);
                    // print('ADDED MSG');
                    return ListView(
                      shrinkWrap: true,
                      reverse: false,
                      padding: EdgeInsets.only(top: 10, bottom: 10),
                      physics: AlwaysScrollableScrollPhysics(),
                      children: messages
                          .map((message) => MessageBubble(
                                message: message,
                                userID: widget.userID,
                              ))
                          .toList(),
                    );
                  },
                ),
              ),
            ),
            Align(
              alignment: Alignment.bottomLeft,
              child: Container(
                height: 84.0,
                padding: EdgeInsets.only(top: 16),
                margin: EdgeInsets.symmetric(horizontal: 16),
                constraints: BoxConstraints(
                  maxWidth: 324,
                ),
                child: TextField(
                  controller: _controller,
                  decoration: InputDecoration(
                    filled: true,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

const MESSAGE_MAX_CHARS = 20;

class MessageBubble extends StatelessWidget {
  MessageBubble({
    Key? key,
    required this.message,
    required this.userID,
  }) : super(key: key);

  final MessageResponse message;
  final Int64 userID;

  @override
  Widget build(BuildContext context) {
    final bool isRight = this.message.user.id == this.userID;

    return Container(
      child: Align(
        alignment: isRight ? Alignment.topRight : Alignment.topLeft,
        child: Container(
          constraints: BoxConstraints(
            minWidth: 184.0,
            maxWidth: 284.0,
          ),
          padding: EdgeInsets.all(16),
          margin: EdgeInsets.symmetric(vertical: 4),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16.0),
            color: isRight
                ? Theme.of(context).primaryColorLight
                : Theme.of(context).primaryColor,
          ),
          child: message.text.length > MESSAGE_MAX_CHARS
              ? Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(message.text),
                    SizedBox(height: 4.0),
                    Align(
                      alignment: Alignment.bottomRight,
                      child: Text(timeago.format(DateTime.now())),
                    ),
                  ],
                )
              : Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(message.text),
                    Align(
                      alignment: Alignment.bottomRight,
                      child: Text(
                          timeago.format(DateTime.parse(message.createdAt))),
                    ),
                  ],
                ),
        ),
      ),
    );
  }
}
