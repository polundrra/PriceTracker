package rabbit

import "github.com/streadway/amqp"

func (r *Rabbit) Publish(data []byte) error {
	_, err := r.Channel.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return r.Channel.Publish(
		"",
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		},
	)
}
