import 'package:flutter/material.dart';

class RegisterScreen extends StatefulWidget {
  RegisterScreen({Key? key}) : super(key: key);

  @override
  _RegisterScreenState createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<RegisterScreen> {
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

  void onRegister() async {
    print('Registering user... ${_controller.value.text}');
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding: const EdgeInsets.only(bottom: 32.0),
              child: FlutterLogo(
                size: 124.0,
              ),
            ),
            Text(
              'Welcome to the chat!'.toUpperCase(),
              style: TextStyle(
                fontSize: 18.0,
                letterSpacing: 2.0,
              ),
            ),
            SizedBox(height: 18.0),
            Container(
              constraints: BoxConstraints(
                maxWidth: 256.0,
              ),
              child: Column(
                children: [
                  TextField(
                    decoration: InputDecoration(
                      labelText: 'Username...',
                    ),
                    controller: _controller,
                  ),
                  SizedBox(height: 32.0),
                  Row(
                    children: [
                      Expanded(
                        child: ElevatedButton(
                          onPressed: onRegister,
                          child: Text('Enter chat'),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
