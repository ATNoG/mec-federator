package config

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

var (
	Producer sarama.SyncProducer
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

	log.Println("Kafka initialized successfully")

	return nil
}
