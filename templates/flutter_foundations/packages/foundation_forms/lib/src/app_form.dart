import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';

class AppForm extends StatelessWidget {
  const AppForm({
    super.key,
    required this.formKey,
    required this.child,
    this.autovalidateMode = AutovalidateMode.disabled,
  });
  final GlobalKey<FormBuilderState> formKey;
  final Widget child;
  final AutovalidateMode autovalidateMode;
  @override
  Widget build(BuildContext context) => FormBuilder(
    key: formKey,
    autovalidateMode: autovalidateMode,
    child: child,
  );
}
