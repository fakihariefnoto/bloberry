import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'api/api_client.dart';
import 'screens/login_screen.dart';
import 'screens/files_screen.dart';
import 'screens/otp_login_screen.dart';

void main() {
  runApp(BloberryApp());
}

class BloberryApp extends StatefulWidget {
  const BloberryApp({super.key});

  @override
  State<BloberryApp> createState() => _BloberryAppState();
}

class _BloberryAppState extends State<BloberryApp> {
  late final ApiClient api;
  late final GoRouter router;

  @override
  void initState() {
    super.initState();
    api = ApiClient(baseUrl: const String.fromEnvironment('BLOBERRY_SERVER', defaultValue: 'http://localhost:8080'));
    router = GoRouter(
      initialLocation: '/',
      routes: [
        GoRoute(path: '/', builder: (c, s) => LoginScreen(api: api)),
        GoRoute(path: '/otp', builder: (c, s) => OtpLoginScreen(api: api)),
        GoRoute(path: '/files', builder: (c, s) => FilesScreen(api: api)),
      ],
    );
    _bootstrap();
  }

  Future<void> _bootstrap() async {
    await api.loadTokens();
    if (mounted && api.accessToken != null) {
      router.go('/files');
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Bloberry',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF8B7DEB)),
        useMaterial3: true,
        fontFamily: null, // system font — SF Pro / Roboto (mobile/README.md)
      ),
      routerConfig: router,
    );
  }
}
