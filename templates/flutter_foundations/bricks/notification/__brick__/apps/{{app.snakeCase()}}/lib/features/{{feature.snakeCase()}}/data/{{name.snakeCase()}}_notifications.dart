import 'package:foundation_notifications/foundation_notifications.dart';

class {{name.pascalCase()}}Notifications {
  const {{name.pascalCase()}}Notifications(this.notifications);
  final LocalNotificationService notifications;
  Future<void> notify(String message) {
    return notifications.show(
      id: message.hashCode,
      title: '{{name.titleCase()}}',
      body: message,
    );
  }
}
