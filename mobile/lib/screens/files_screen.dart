import 'package:dio/dio.dart';
import 'package:flutter/material.dart';

import '../api/api_client.dart';
import 'shares_screen.dart';
import 'uploads_screen.dart';
import 'profile_screen.dart';

class FilesScreen extends StatefulWidget {
  const FilesScreen({super.key, required this.api});
  final ApiClient api;

  @override
  State<FilesScreen> createState() => _FilesScreenState();
}

class _FilesScreenState extends State<FilesScreen> {
  int _tab = 0;
  List<Map<String, dynamic>> _folders = [];
  List<Map<String, dynamic>> _objects = [];
  bool _loading = true;
  String? _error;
  String _currentFolder = 'root';

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final res = await widget.api.dio.get('/folders/$_currentFolder/children');
      final data = res.data['data'] as Map;
      setState(() {
        _folders = (data['folders'] as List? ?? []).cast<Map<String, dynamic>>();
        _objects = (data['objects'] as List? ?? []).cast<Map<String, dynamic>>();
      });
    } on DioException catch (e) {
      final apiErr = e.error;
      setState(() => _error = apiErr is ApiException ? apiErr.content : 'Could not load files');
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<void> _createFolder() async {
    final name = await _prompt('New folder', 'Folder name');
    if (name == null || name.isEmpty) return;
    try {
      await widget.api.dio.post('/folders', data: {'name': name, 'parent_id': _currentFolder});
      _load();
    } catch (_) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Could not create folder')));
    }
  }

  Future<String?> _prompt(String title, String label) {
    final controller = TextEditingController();
    return showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(title),
        content: TextField(controller: controller, decoration: InputDecoration(labelText: label), autofocus: true),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancel')),
          FilledButton(onPressed: () => Navigator.pop(ctx, controller.text.trim()), child: const Text('Create')),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _tab,
        children: [
          _buildFiles(context),
          const UploadsScreen(),
          const SharesScreen(),
          const ProfileScreen(),
        ],      ),
      bottomNavigationBar: NavigationBar(
        selectedIndex: _tab,
        onDestinationSelected: (i) => setState(() => _tab = i),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.folder_outlined), selectedIcon: Icon(Icons.folder), label: 'Files'),
          NavigationDestination(icon: Icon(Icons.cloud_upload_outlined), selectedIcon: Icon(Icons.cloud_upload), label: 'Uploads'),
          NavigationDestination(icon: Icon(Icons.share_outlined), selectedIcon: Icon(Icons.share), label: 'Shares'),
          NavigationDestination(icon: Icon(Icons.person_outline), selectedIcon: Icon(Icons.person), label: 'More'),
        ],
      ),
    );
  }

  Widget _buildFiles(BuildContext context) {
    return SafeArea(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                const Text('Files', style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
                const Spacer(),
                IconButton(
                  tooltip: 'New folder',
                  onPressed: _createFolder,
                  icon: const Icon(Icons.create_new_folder_outlined),
                ),
                IconButton(tooltip: 'Refresh', onPressed: _load, icon: const Icon(Icons.refresh)),
              ],
            ),
          ),
          Expanded(child: _buildList()),
        ],
      ),
    );
  }

  Widget _buildList() {
    if (_loading) {
      return ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 16),
        itemCount: 8,
        itemBuilder: (_, i) => Container(
          height: 44,
          margin: const EdgeInsets.only(bottom: 8),
          decoration: BoxDecoration(color: Colors.grey.shade200, borderRadius: BorderRadius.circular(8)),
        ),
      );
    }
    if (_error != null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(_error!, style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 8),
            FilledButton.tonal(onPressed: _load, child: const Text('Retry')),
          ],
        ),
      );
    }
    if (_folders.isEmpty && _objects.isEmpty) {
      return const Center(child: Text('No files here yet'));
    }
    return ListView(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      children: [
        for (final f in _folders)
          ListTile(
            leading: const Icon(Icons.folder, color: Colors.indigo),
            title: Text(f['name'] as String? ?? ''),
            onTap: () {
              setState(() => _currentFolder = f['id'] as String);
              _load();
            },
          ),
        for (final o in _objects)
          ListTile(
            leading: Icon(o['visibility'] == 'public' ? Icons.public : Icons.insert_drive_file_outlined,
                color: o['visibility'] == 'public' ? Colors.orange : Colors.grey),
            title: Text(o['name'] as String? ?? ''),
            trailing: o['visibility'] == 'public'
                ? const Text('public', style: TextStyle(color: Colors.orange, fontSize: 12))
                : null,
            onTap: () => _download(o),
          ),
      ],
    );
  }

  Future<void> _download(Map<String, dynamic> obj) async {
    try {
      final id = obj['id'] as String;
      final name = obj['name'] as String? ?? 'download';
      final res = await widget.api.dio.get('/objects/$id/download');
      if (res.statusCode == 302 || (res.data is Map && res.data['data'] != null)) {
        final data = res.data is Map ? (res.data['data'] as Map?) : null;
        final url = data?['download_url'] ?? data?['url'];
        if (url != null && mounted) {
          ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Opening: $name')));
        }
      }
    } catch (_) {
      // v1: streaming download on mobile is a follow-up; surface a snackbar.
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Download not available on mobile yet')));
      }
    }
  }
}
