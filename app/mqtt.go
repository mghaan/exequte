package app

import (
	"fmt"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mghaan/exequte/logger"
)

type Server struct {
	broker *paho.ClientOptions
	client paho.Client
	logs   *logger.Logger
}

func StartMqtt(logger *logger.Logger, cfg *Config) *Server {
	server := &Server{broker: paho.NewClientOptions(), logs: logger}

	server.Connect(cfg.Mqtt.Ssl, cfg.Mqtt.Host, cfg.Mqtt.Port, cfg.Mqtt.Client, cfg.Mqtt.User, cfg.Mqtt.Password)

	return server
}

// Log event.
func (server *Server) Log() *logger.Logger {
	return server.logs
}

// Connect to MQTT server.
func (server *Server) Connect(ssl bool, host string, port int, clientid string, username string, password string) {
	proto := "tcp"
	if ssl {
		proto = "ssl"
	}
	conn := proto + "://" + host + ":" + fmt.Sprintf("%d", port)
	server.broker.AddBroker(conn)
	server.broker.SetClientID(clientid)
	server.broker.SetUsername(username)
	server.broker.SetPassword(password)
	server.broker.SetWill(clientid+"/online", "0", 0, false)
	server.broker.SetConnectionLostHandler(server.handlerConnectionLost)

	server.client = paho.NewClient(server.broker)

	server.handlerConnectionAttempt()
}

func (server *Server) handlerConnectionLost(client paho.Client, err error) {
	server.logs.Error(logger.MQTT, "Connection lost", err)
	server.handlerConnectionAttempt()
}

func (server *Server) handlerConnectionAttempt() {
	var err error

	at := 0
	for {
		at++

		token := server.client.Connect()
		token.Wait()
		if err = token.Error(); err == nil {
			server.logs.Info(logger.MQTT, fmt.Sprintf("Connected to %s:%s", server.broker.Servers[0].Hostname(), server.broker.Servers[0].Port()))
			server.Publish("online", "1")

			return
		}

		if at == 3 {
			server.logs.Fatal(logger.MQTT, fmt.Sprintf("Connection to %s:%s failed", server.broker.Servers[0].Hostname(), server.broker.Servers[0].Port()), err)
		}

		time.Sleep(30 * time.Second)
	}
}

// Subscribe to this topic.
func (server *Server) Subscribe(topic string, callback paho.MessageHandler) bool {
	topic = server.broker.ClientID + "/" + topic
	token := server.client.Subscribe(topic, 0, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		server.logs.Error(logger.MQTT, fmt.Sprintf("Unable to subscribe topic '%s'", topic), err)

		return false
	}

	server.logs.Info(logger.MQTT, fmt.Sprintf("Subscribed topic '%s'", topic))

	return true
}

// Disconnect from server.
func (server *Server) Disconnect() {
	if server.client != nil {
		server.Publish("online", "0")
		server.client.Disconnect(0)
		server.logs.Info(logger.MQTT, "Disconnected")
	}
}

// Send value to this topic.
func (server *Server) Publish(topic string, value string) bool {
	topic = server.broker.ClientID + "/" + topic

	token := server.client.Publish(topic, 0, false, value)
	token.Wait()
	if err := token.Error(); err != nil {
		server.logs.Info(logger.MQTT, fmt.Sprintf("Failed to publish '%s' -> %s", topic, value))

		return false
	}

	server.logs.Info(logger.MQTT, fmt.Sprintf("Publish '%s' -> %s", topic, value))

	return true
}
