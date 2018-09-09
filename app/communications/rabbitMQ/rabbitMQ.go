package rabbitMQ

import (
	"github.com/streadway/amqp"
	"theAmazingPostManager/app/common"
)

type RabbitMQTask interface {
	GetMessageBytes() ([]byte, error)
	GetQueue() string
}

func PublishMessageOnExchange(newTask RabbitMQTask,exchangeName string,routingKey string) error{

	ch := common.GetRabbitMQChannel()
	defer ch.Close()

	err := ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messageBody, err := newTask.GetMessageBytes()
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchangeName, // exchange
		routingKey,     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        messageBody,
		})

	if err != nil {
		return err
	}

	return nil

}

