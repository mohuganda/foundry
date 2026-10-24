import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:foundation_forms/foundation_forms.dart';

class {{name.pascalCase()}}Form extends StatefulWidget {
  const {{name.pascalCase()}}Form({super.key});
  @override
  State<{{name.pascalCase()}}Form> createState() => _{{name.pascalCase()}}FormState();
}

class _{{name.pascalCase()}}FormState extends State<{{name.pascalCase()}}Form> {
  final formKey = GlobalKey<FormBuilderState>();
  @override
  Widget build(BuildContext context) => AppForm(
    formKey: formKey,
    child: Column(
      children: [
        AppTextField(
          name: 'name',
          label: 'Name',
          validator: FormBuilderValidators.required(),
        ),
        const SizedBox(height: 16),
        FilledButton(
          onPressed: () {
            if (formKey.currentState?.saveAndValidate() ?? false) {}
          },
          child: const Text('Save'),
        ),
      ],
    ),
  );
}
