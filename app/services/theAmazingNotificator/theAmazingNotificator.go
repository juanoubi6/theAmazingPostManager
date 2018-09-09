package theAmazingNotificator

import (
	"theAmazingPostManager/app/communications/rabbitMQ"
	"theAmazingPostManager/app/config"
)

func SendNotification(notificationTask rabbitMQ.RabbitMQTask,routingKey string) error{

	if err := rabbitMQ.PublishMessageOnExchange(notificationTask,config.GetConfig().RABBIT_NOTIFICATION_EXCHANGE,routingKey); err != nil {
		return err
	}

	return nil

}