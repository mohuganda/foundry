import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';

class AppTextField extends StatelessWidget {
  const AppTextField({
    super.key,
    required this.name,
    this.label,
    this.validator,
    this.maxLines = 1,
  });
  final String name;
  final String? label;
  final String? Function(String?)? validator;
  final int maxLines;
  @override
  Widget build(BuildContext context) => FormBuilderTextField(
    name: name,
    validator: validator,
    maxLines: maxLines,
    decoration: InputDecoration(
      labelText: label,
      border: const OutlineInputBorder(),
    ),
  );
}

class AppDateField extends StatelessWidget {
  const AppDateField({super.key, required this.name, this.label});
  final String name;
  final String? label;
  @override
  Widget build(BuildContext context) => FormBuilderDateTimePicker(
    name: name,
    inputType: InputType.date,
    decoration: InputDecoration(
      labelText: label,
      border: const OutlineInputBorder(),
    ),
  );
}
