# Sistema de Notificações Assíncronas

Sistema de notificações desenvolvido com Flutter (frontend) e Go (backend), utilizando RabbitMQ para processamento assíncrono de mensagens.

## Tecnologias Utilizadas

Backend: Go + Gin Framework + RabbitMQ
Frontend: Flutter + HTTP Package
Mensageria: RabbitMQ (CloudAMQP)

## Como Executar

Pré-requisitos

Go 1.19+ instalado
Flutter 3.0+ instalado

1. Backend (Go)
   bashcd backend
   go mod tidy
   go run main.go
   O servidor estará disponível em http://localhost:8080
2. Frontend (Flutter)
   bashcd frontend
   flutter pub get
   flutter run

## API Endpoints

POST /api/notificar - Enviar notificação
GET /api/notificacao/status/{id} - Consultar status

## Fluxo da Aplicação

Flutter envia mensagem para o backend
Backend publica mensagem no RabbitMQ
Consumidor processa mensagem assincronamente
Frontend consulta status periodicamente
Status é atualizado em tempo real na interface
