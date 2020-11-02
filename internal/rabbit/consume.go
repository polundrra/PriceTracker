package rabbit

import (
	"github.com/polundrra/PriceTracker/internal/mail"
	"log"
)

func (r *Rabbit) StartConsumer() error {
	_, err := r.Channel.QueueDeclare(
		r.QueueName,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	prefetchCount := r.Concurrency * 4
	err = r.Channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return err
	}

	msgs, err := r.Channel.Consume(
		r.QueueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	for i := 0; i < r.Concurrency; i++ {
		log.Printf("Processing messages on thread %v...\n", i)
		go func() {
			for msg := range msgs {
				err := mail.SendEmail(msg.Body)
				if err == nil {
					if err = msg.Ack(false); err != nil {
						log.Printf("error acknowledge message: %v", err)
					}
				} else {
					if err = msg.Nack(false, true); err != nil {
						log.Printf("error nacknoledge message: %v", err)
					}
				}
			}
		}()
	}
	return nil
}
