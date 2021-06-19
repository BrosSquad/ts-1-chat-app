import 'package:flutter/material.dart';

import './screens/register_screen.dart';

void main() {
  runApp(ChatApp());
}

class ChatApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Chat',
      theme: ThemeData.dark(),
      home: RegisterScreen(),
    );
  }
}
