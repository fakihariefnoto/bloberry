import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../api/api_client.dart';

class OtpLoginScreen extends StatefulWidget {
  const OtpLoginScreen({super.key, required this.api});
  final ApiClient api;

  @override
  State<OtpLoginScreen> createState() => _OtpLoginScreenState();
}

class _OtpLoginScreenState extends State<OtpLoginScreen> {
  final _email = TextEditingController();
  final _code = TextEditingController();
  bool _sent = false;
  bool _loading = false;
  String? _error;

  Future<void> _request() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      await widget.api.dio.post('/auth/otp/request', data: {'email': _email.text.trim()});
      setState(() => _sent = true);
    } on DioException catch (e) {
      final apiErr = e.error;
      setState(() => _error = apiErr is ApiException ? apiErr.content : 'Could not reach the server');
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<void> _verify() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final res = await widget.api.dio.post('/auth/otp/verify', data: {
        'email': _email.text.trim(),
        'code': _code.text.trim(),
        'platform': 'mobile',
      });
      final data = res.data['data'] as Map;
      await widget.api.saveTokens(
        access: data['access_token'] as String,
        refresh: data['refresh_token'] as String,
      );
      if (mounted) context.go('/files');
    } on DioException catch (e) {
      final apiErr = e.error;
      setState(() => _error = apiErr is ApiException ? apiErr.content : 'Could not reach the server');
    } finally {
      setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Email login code')),
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Text(_sent ? 'Enter the 6-digit code we emailed you.' : 'Enter your email to receive a code.'),
                  const SizedBox(height: 16),
                  if (!_sent)
                    TextField(
                      controller: _email,
                      decoration: const InputDecoration(labelText: 'Email', border: OutlineInputBorder()),
                      keyboardType: TextInputType.emailAddress,
                    )
                  else
                    TextField(
                      controller: _code,
                      decoration: const InputDecoration(labelText: '6-digit code', border: OutlineInputBorder()),
                      keyboardType: TextInputType.number,
                    ),
                  if (_error != null) ...[
                    const SizedBox(height: 12),
                    Text(_error!, style: const TextStyle(color: Colors.red)),
                  ],
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: _loading ? null : (_sent ? _verify : _request),
                    style: FilledButton.styleFrom(minimumSize: const Size.fromHeight(48)),
                    child: _loading
                        ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2))
                        : Text(_sent ? 'Verify' : 'Send code'),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
