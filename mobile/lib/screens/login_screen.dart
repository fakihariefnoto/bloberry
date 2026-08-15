import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../api/api_client.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key, required this.api});
  final ApiClient api;

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _email = TextEditingController();
  final _password = TextEditingController();
  bool _loading = false;
  String? _error;

  Future<void> _login() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final res = await widget.api.dio.post('/auth/login', data: {
        'email': _email.text.trim(),
        'password': _password.text,
        'platform': 'mobile',
      });
      final data = res.data['data'] as Map;
      final totp = data['totp_required'] == true;
      await widget.api.saveTokens(
        access: data['access_token'] as String,
        refresh: data['refresh_token'] as String,
      );
      if (!mounted) return;
      if (totp) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Enter your 2FA code to finish signing in')),
        );
        // v1: TOTP step on mobile is a follow-up; land on files for now.
      }
      context.go('/files');
    } on DioException catch (e) {
      final apiErr = e.error;
      if (apiErr is ApiException) {
        setState(() => _error = apiErr.content);
      } else {
        setState(() => _error = 'Could not reach the server');
      }
    } catch (_) {
      setState(() => _error = 'Could not reach the server');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const Text('Bloberry', style: TextStyle(fontSize: 28, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  const Text('Sign in to your storage', style: TextStyle(color: Colors.grey)),
                  const SizedBox(height: 24),
                  TextField(
                    controller: _email,
                    decoration: const InputDecoration(labelText: 'Email', border: OutlineInputBorder()),
                    keyboardType: TextInputType.emailAddress,
                  ),
                  const SizedBox(height: 12),
                  TextField(
                    controller: _password,
                    obscureText: true,
                    decoration: const InputDecoration(labelText: 'Password', border: OutlineInputBorder()),
                  ),
                  if (_error != null) ...[
                    const SizedBox(height: 12),
                    Text(_error!, style: const TextStyle(color: Colors.red)),
                  ],
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: _loading ? null : _login,
                    style: FilledButton.styleFrom(minimumSize: const Size.fromHeight(48)),
                    child: _loading
                        ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                        : const Text('Log in'),
                  ),
                  const SizedBox(height: 12),
                  TextButton(onPressed: () => context.go('/otp'), child: const Text('Use a code instead')),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
