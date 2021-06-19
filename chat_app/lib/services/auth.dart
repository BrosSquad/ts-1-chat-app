import 'package:chat_app/proto/auth.pbgrpc.dart';
import 'package:fixnum/fixnum.dart';
import 'package:grpc/grpc.dart';

class AuthService {
  late AuthClient client;
  late Int64 userID;
  late String username;

  AuthService(ClientChannel channel) {
    client = AuthClient(channel);
  }

  Future<void> register(String username) async {
    try {
      RegisterResponse res =
          await client.register(RegisterRequest(username: username));

      this.username = res.user.username;
      this.userID = res.user.id;

      print(
          'Successfully registered user with username ${this.username} and ID ${this.userID}');
    } catch (error) {
      print('Error registering user $error');
    }
  }
}
