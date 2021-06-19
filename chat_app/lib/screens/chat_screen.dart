import 'package:chat_app/proto/chat.pb.dart';
import 'package:chat_app/proto/user.pb.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:timeago/timeago.dart' as timeago;

import '../grpc_client.dart';

User userA = User(username: 'Misko');
User userB = User(username: 'Pisko');

List<MessageResponse> messages = [
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text:
        'eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Eo buraz',
    user: userB,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Pisem jebozovni stack',
    user: userB,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
  MessageResponse(
    createdAt: 'Sat, 19 Jun 2021 20:16:24 +0200',
    text: 'Ne seeeeri',
    user: userA,
  ),
];

class ChatScreen extends StatefulWidget {
  const ChatScreen({Key? key}) : super(key: key);

  @override
  _ChatScreenState createState() => _ChatScreenState();
}

class _ChatScreenState extends State<ChatScreen> {
  late TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void sendMessage() {}

  @override
  Widget build(BuildContext context) {
    String username = GRPCClient.of(context).authService.username;
    Int64 userID = GRPCClient.of(context).authService.userID;

    return Scaffold(
      appBar: AppBar(
        title: Text(username),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: sendMessage,
        child: Icon(Icons.send),
      ),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: ListView.builder(
                  itemCount: messages.length,
                  shrinkWrap: true,
                  reverse: true,
                  padding: EdgeInsets.only(top: 10, bottom: 10),
                  physics: AlwaysScrollableScrollPhysics(),
                  itemBuilder: (context, index) =>
                      MessageBubble(message: messages[index]),
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

class MessageBubble extends StatelessWidget {
  const MessageBubble({
    Key? key,
    required this.message,
  }) : super(key: key);

  final MessageResponse message;

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Align(
        alignment: message.user.username == 'Misko'
            ? Alignment.topRight
            : Alignment.topLeft,
        child: Container(
          constraints: BoxConstraints(
            minWidth: 184.0,
            maxWidth: 284.0,
          ),
          padding: EdgeInsets.all(16),
          margin: EdgeInsets.symmetric(vertical: 4),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16.0),
            color: message.user.username == 'Misko'
                ? Theme.of(context).primaryColorLight
                : Theme.of(context).primaryColor,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(message.text),
              SizedBox(height: 4.0),
              Align(
                alignment: Alignment.bottomRight,
                child: Text(timeago.format(DateTime.now())),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
