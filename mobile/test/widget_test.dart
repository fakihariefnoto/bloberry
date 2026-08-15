import 'package:flutter_test/flutter_test.dart';

import 'package:bloberry/main.dart';

void main() {
  testWidgets('app boots to the login screen', (WidgetTester tester) async {
    await tester.pumpWidget(const BloberryApp());
    await tester.pump();

    expect(find.text('Bloberry'), findsOneWidget);
    expect(find.text('Log in'), findsOneWidget);
  });
}
