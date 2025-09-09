import 'package:flutter/material.dart';

class BuildStatus extends StatelessWidget {
  final String status;
  const BuildStatus({super.key, required this.status});

  @override
  Widget build(BuildContext context) {
    Color cor;
    String texto;

    if (status == 'PROCESSADO_SUCESSO') {
      cor = Colors.green;
      texto = 'Processado';
    } else if (status.contains('AGUARDANDO') || status.contains('RECEBIDA')) {
      cor = Colors.orange;
      texto = 'Aguardando';
    } else {
      cor = Colors.red;
      texto = 'Erro';
    }

    return Container(
      padding: EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: cor,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(
        texto,
        style: TextStyle(color: Colors.white, fontSize: 12),
      ),
    );
  }
}
