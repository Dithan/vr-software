package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Dados que vão e voltam
type Mensagem struct {
	MensagemID       string `json:"mensagemId"`
	ConteudoMensagem string `json:"conteudoMensagem"`
}

// Armazenar status em memória
var statusMensagens = make(map[string]string)

// Conexão RabbitMQ
var canal *amqp.Channel

func main() {
	// 1. Conectar RabbitMQ
	conectarRabbitMQ()

	// 2. Criar filas
	criarFilas()

	// 3. Começar a ouvir mensagens
	go processarMensagens()

	// 4. Criar servidor web
	servidor := gin.Default()
	servidor.Use(cors.Default()) // Permite Flutter chamar

	// 5. Criar rotas
	servidor.POST("/api/notificar", enviarNotificacao)
	servidor.GET("/api/notificacao/status/:id", consultarStatus)

	log.Println("🚀 Servidor rodando em http://localhost:8080")
	servidor.Run(":8080")
}

func conectarRabbitMQ() {
	conn, err := amqp.Dial("amqp://bjnuffmq:gj-YQIiEXyfxQxjsZtiYDKeXIT8ppUq7@jaragua-01.lmq.cloudamqp.com/bjnuffmq")
	if err != nil {
		log.Fatal("❌ Não conseguiu conectar RabbitMQ:", err)
	}

	canal, err = conn.Channel()
	if err != nil {
		log.Fatal("❌ Não conseguiu abrir canal:", err)
	}

	log.Println("✅ RabbitMQ conectado!")
}

func criarFilas() {
	// Fila onde chegam as mensagens
	filaEntrada := "fila.notificacao.entrada.thiago"
	_, err := canal.QueueDeclare(filaEntrada, true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Erro na fila entrada:", err)
	}

	// Fila onde vão os status
	filaStatus := "fila.notificacao.status.thiago"
	_, err = canal.QueueDeclare(filaStatus, true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Erro na fila status:", err)
	}

	log.Println("✅ Filas criadas!")
}

// Recebe mensagem do Flutter
func enviarNotificacao(c *gin.Context) {
	var msg Mensagem

	// Pegar dados do Flutter
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(400, gin.H{"erro": "Dados inválidos"})
		return
	}

	// Criar ID se não veio
	if msg.MensagemID == "" {
		msg.MensagemID = uuid.New().String()
	}

	// Salvar status inicial
	statusMensagens[msg.MensagemID] = "AGUARDANDO_PROCESSAMENTO"

	// Enviar para RabbitMQ
	dadosJson, _ := json.Marshal(msg)
	filaEntrada := "fila.notificacao.entrada.thiago"
	
	err := canal.Publish("", filaEntrada, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        dadosJson,
	})

	if err != nil {
		c.JSON(500, gin.H{"erro": "Falha ao enviar"})
		return
	}

	log.Printf("📨 Mensagem enviada: %s", msg.MensagemID)

	// Responde pro Flutter
	c.JSON(202, gin.H{
		"mensagemId": msg.MensagemID,
		"status":     "RECEBIDA_PARA_PROCESSAMENTO",
	})
}

// Flutter consulta status
func consultarStatus(c *gin.Context) {
	id := c.Param("id")
	status, existe := statusMensagens[id]

	if !existe {
		c.JSON(404, gin.H{"erro": "Mensagem não encontrada"})
		return
	}

	c.JSON(200, gin.H{
		"mensagemId": id,
		"status":     status,
	})
}

// Ouvindo e processando mensagens
func processarMensagens() {
	filaEntrada := "fila.notificacao.entrada.thiago"
	mensagens, err := canal.Consume(filaEntrada, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Erro ao ouvir mensagens:", err)
	}

	log.Println("🔄 Ouvindo mensagens...")

	for mensagem := range mensagens {
		var msg Mensagem
		json.Unmarshal(mensagem.Body, &msg)

		log.Printf("⚙️ Processando: %s", msg.MensagemID)

		// Simula processamento
		time.Sleep(2 * time.Second)

		// Atualiza status
		statusMensagens[msg.MensagemID] = "PROCESSADO_SUCESSO"

		// Envia para fila de status
		filaStatus := "fila.notificacao.status.thiago"
		statusJson, _ := json.Marshal(map[string]string{
			"mensagemId": msg.MensagemID,
			"status":     "PROCESSADO_SUCESSO",
		})

		canal.Publish("", filaStatus, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        statusJson,
		})

		log.Printf("✅ Processado: %s", msg.MensagemID)
	}
}