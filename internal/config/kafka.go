package config

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

var (
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
)

func InitKafka() error {
	brokers := []string{AppConfig.KafkaHost + ":" + AppConfig.KafkaPort}
	log.Println("Brokers:", brokers)

	username := AppConfig.KafkaUsername
	password := AppConfig.KafkaPassword

	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Retry.Backoff = 100 * time.Millisecond

	producerConfig.Net.SASL.Enable = true
	producerConfig.Net.SASL.User = username
	producerConfig.Net.SASL.Password = password
	log.Println("Connecting to Kafka with username:", username, "and password:", password)

	p, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		return err
	}
	Producer = p

	consumerConfig := sarama.NewConfig()
	consumerConfig.Net.SASL.Enable = true
	consumerConfig.Net.SASL.User = username
	consumerConfig.Net.SASL.Password = password

	c, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		return err
	}
	Consumer = c

	log.Println("Kafka initialized successfully")

	return nil
}
