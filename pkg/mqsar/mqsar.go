package mqsar

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/fhmq/hmq/broker"
	"os"
	"os/signal"
)

type Config struct {
	PulsarConfig PulsarConfig
}

type PulsarConfig struct {
	Host     string
	HttpPort int
	TcpPort  int
}

func Run(config *Config, impl Server) (err error) {
	mqttConfig := &broker.Config{}
	mqttConfig.Port = "1883"
	clientOptions := pulsar.ClientOptions{}
	clientOptions.URL = fmt.Sprintf("pulsar://%s:%d", config.PulsarConfig.Host, config.PulsarConfig.TcpPort)
	mqttConfig.Plugin.Bridge, err = newPulsarBridgeMq(clientOptions, impl)
	if err != nil {
		return
	}
	newBroker, err := broker.NewBroker(mqttConfig)
	if err != nil {
		return err
	}
	newBroker.Start()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			return nil
		}
	}
}
