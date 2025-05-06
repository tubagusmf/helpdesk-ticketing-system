package config

import amqp "github.com/rabbitmq/amqp091-go"

func InitRabbitMQ() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/helpdesk_ticketing")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"notification", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"emailQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		"emailQueue",   // queue name
		"emailQueue",   // routing key
		"notification", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
