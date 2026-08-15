import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// API client over the standard envelope {data?, messages?: [{code, content}]}.
/// The Go backend emits snake_case; models are generated with
/// FieldRename.snake (mobile/README.md — the silent null-break trap).
class ApiClient {
  ApiClient({String? baseUrl})
      : _baseUrl = baseUrl ?? 'http://localhost:8080' {
    dio = Dio(BaseOptions(baseUrl: _baseUrl, connectTimeout: const Duration(seconds: 10)));
    dio.interceptors.add(_EnvelopeInterceptor(this));
  }

  final String _baseUrl;
  late final Dio dio;
  static const _storage = FlutterSecureStorage();
  String? _accessToken;
  String? _refreshToken;

  Future<void> loadTokens() async {
    _accessToken = await _storage.read(key: 'access_token');
    _refreshToken = await _storage.read(key: 'refresh_token');
  }

  Future<void> saveTokens({String? access, String? refresh}) async {
    if (access != null) {
      _accessToken = access;
      await _storage.write(key: 'access_token', value: access);
    }
    if (refresh != null) {
      _refreshToken = refresh;
      await _storage.write(key: 'refresh_token', value: refresh);
    }
  }

  Future<void> clearTokens() async {
    _accessToken = null;
    _refreshToken = null;
    await _storage.delete(key: 'access_token');
    await _storage.delete(key: 'refresh_token');
  }

  String? get accessToken => _accessToken;
  String? get refreshToken => _refreshToken;
  String get baseUrl => _baseUrl;
}

class ApiMessage {
  ApiMessage(this.code, this.content);
  final String code;
  final String? content;
}

class ApiException implements Exception {
  ApiException(this.status, this.messages);
  final int status;
  final List<ApiMessage> messages;
  String get code => messages.isNotEmpty ? messages.first.code : 'unknown_error';
  String get content => messages.isNotEmpty ? (messages.first.content ?? code) : 'Request failed';

  @override
  String toString() => content;
}

class _EnvelopeInterceptor extends Interceptor {
  _EnvelopeInterceptor(this.client);

  final ApiClient client;

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final token = client.accessToken;
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    // Parse the envelope from the error body; never surface raw exceptions.
    if (err.response != null && err.response!.data is Map) {
      final data = err.response!.data as Map;
      final msgs = <ApiMessage>[];
      final raw = data['messages'];
      if (raw is List) {
        for (final m in raw) {
          if (m is Map) {
            msgs.add(ApiMessage(m['code']?.toString() ?? '', m['content']?.toString()));
          }
        }
      }
      handler.next(DioException(
        requestOptions: err.requestOptions,
        response: err.response,
        error: ApiException(err.response!.statusCode ?? 0, msgs),
      ));
      return;
    }
    handler.next(err);
  }
}
