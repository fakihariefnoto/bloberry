import 'package:flutter/material.dart';

/// Uploads tab — the persistent queue (mobile/README.md §The upload queue).
/// v1 keeps the queue model + UI; multipart resume is wired in the queue store.
class UploadsScreen extends StatelessWidget {
  const UploadsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Padding(
            padding: EdgeInsets.all(16),
            child: Text('Uploads', style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold)),
          ),
          Expanded(
            child: Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.cloud_upload_outlined, size: 48, color: Colors.grey.shade400),
                  const SizedBox(height: 12),
                  const Text('No uploads in progress'),
                  const SizedBox(height: 4),
                  Text('Uploads resume after a network drop or app restart.',
                      style: TextStyle(color: Colors.grey.shade600, fontSize: 12)),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
