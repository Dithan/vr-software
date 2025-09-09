import 'dart:async';
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:frontend/models/mensagem.dart';
import 'package:frontend/widgets/build_status.dart';
import 'package:uuid/uuid.dart';
import 'package:http/http.dart' as http;

class NotificationScreen extends StatefulWidget {
  const NotificationScreen({super.key});

  @override
  State<NotificationScreen> createState() => _NotificationScreenState();
}

class _NotificationScreenState extends State<NotificationScreen> {
  final _controller = TextEditingController();
  final List<Mensagem> _mensagens = [];
  Timer? _timer;

  // URL Backend
  final String url = 'http://localhost:8080';

  @override
  void initState() {
    super.initState();
    _timer = Timer.periodic(Duration(seconds: 3), (_) => _verificaStatus());
  }

  @override
  void dispose() {
    _timer?.cancel();
    _controller.dispose();
    super.dispose();
  }

  // Função que envia a mensagem
  Future<void> _enviarMensagem() async {
    if (_controller.text.isEmpty) return;

    final id = Uuid().v4();
    final texto = _controller.text;

    setState(() {
      _mensagens.insert(0, Mensagem(id, texto, 'AGUARDANDO_PROCESSAMENTO'));
    });

    try {
      // Envia no Backend
      final response = await http.post(
        Uri.parse('$url/api/notificar'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'mensagemId': id,
          'conteudoMensagem': texto,
        }),
      );

      if (response.statusCode == 202) {
        _controller.clear();
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('✅ Mensagem enviada!'),
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('❌ Erro na conexão'),
        ),
      );
    }
  }

  // Verifica o status das mensagens
  Future<void> _verificaStatus() async {
    for (var msg in _mensagens) {
      if (msg.status == 'AGUARDANDO_PROCESSAMENTO' ||
          msg.status == 'RECEBIDA_PARA_PROCESSAMENTO') {
        try {
          final response = await http.get(
            Uri.parse('$url/api/notificacao/status/${msg.id}'),
          );
          if (response.statusCode == 200) {
            final data = jsonDecode(response.body);
            setState(() {
              msg.status = data['status'];
            });
          }
        } catch (e) {
          print('Erro: $e');
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('VR Soft Notificações'),
        backgroundColor: Colors.blue,
      ),
      body: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          children: [
            // Campo de texto
            TextField(
              controller: _controller,
              decoration: InputDecoration(
                labelText: 'Digite sua mensagem',
                border: OutlineInputBorder(),
              ),
            ),

            SizedBox(height: 16),

            // Botão enviar
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _enviarMensagem,
                child: Text('Enviar Notificação'),
              ),
            ),

            SizedBox(height: 16),

            // Lista de mensagens
            Expanded(
              child: _mensagens.isEmpty
                  ? Center(child: Text('Nenhuma mensagem enviada'))
                  : ListView.builder(
                      itemCount: _mensagens.length,
                      itemBuilder: (context, index) {
                        final msg = _mensagens[index];
                        return Card(
                          child: ListTile(
                            title: Text(msg.texto),
                            subtitle: Text('ID: ${msg.id.substring(0, 8)}...'),
                            trailing: BuildStatus(status: msg.status),
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
