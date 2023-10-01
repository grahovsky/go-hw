package rmq

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       *amqp.Queue
	contentType string
}

func (q *Queue) Close() error {
	if err := q.channel.Close(); err != nil {
		logger.Error(fmt.Sprintf("failed to close RMQ channel: %v", err))
	}
	return q.conn.Close()
}

func (q *Queue) Push(ctx context.Context, msg []byte) error {
	return q.channel.PublishWithContext(
		ctx,
		"",           // exchange
		q.queue.Name, // routing
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  q.contentType,
			Body:         msg,
			Priority:     0,
		})
}

func (q *Queue) ConsumeChannel(ctx context.Context, consumer string) (<-chan []byte, error) {
	delivery, err := q.channel.Consume(q.queue.Name, consumer, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("consume channel: %w", err)
	}

	ch := make(chan []byte)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-delivery:
				if !ok {
					return
				}
				ch <- d.Body
			}
		}
	}()
	return ch, err
}

func NewQueue(rmqCf *config.RMQ, contentType string) (*Queue, error) {
	conn, err := amqp.Dial(amqURI(rmqCf))
	if err != nil {
		return nil, fmt.Errorf("connect to rmq %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %w", err)
	}

	queue, err := channel.QueueDeclare(rmqCf.Queue, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("queue declare: %w", err)
	}

	return &Queue{
		conn:        conn,
		channel:     channel,
		queue:       &queue,
		contentType: contentType,
	}, nil
}

func amqURI(rmqCf *config.RMQ) string {
	port, err := strconv.Atoi(rmqCf.Port)
	if err != nil {
		logger.Error(err.Error())
		port = 5672
	}

	return amqp.URI{
		Scheme:   "amqp",
		Host:     rmqCf.Host,
		Port:     port,
		Username: rmqCf.User,
		Password: rmqCf.Password,
		Vhost:    "/",
	}.String()
}
