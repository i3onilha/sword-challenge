package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"sword-challenge/internal/models"
	"sword-challenge/internal/repository"
)

type NotificationConsumer struct {
	broker           MessageBroker
	userRepo         repository.UserRepository
	notificationRepo repository.NotificationRepository
}

func NewNotificationConsumer(
	broker MessageBroker,
	userRepo repository.UserRepository,
	notificationRepo repository.NotificationRepository,
) *NotificationConsumer {
	return &NotificationConsumer{
		broker:           broker,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	// Type assert to get the RabbitMQ implementation for consuming messages
	rabbitmq, ok := c.broker.(*RabbitMQ)
	if !ok {
		return fmt.Errorf("broker is not a RabbitMQ implementation")
	}

	msgs, err := rabbitmq.channel.Consume(
		TaskCreatedQueue, // queue
		"",               // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			ctx := context.Background()
			var taskMsg struct {
				TaskID       int64  `json:"task_id"`
				TechnicianID int64  `json:"technician_id"`
				Title        string `json:"title"`
			}

			if err := json.Unmarshal(msg.Body, &taskMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				msg.Nack(false, false)
				continue
			}

			// Get technician details
			technician, err := c.userRepo.GetByID(ctx, taskMsg.TechnicianID)
			if err != nil {
				log.Printf("Error getting technician: %v", err)
				msg.Nack(false, false)
				continue
			}

			// Create notification
			task := &models.Task{
				ID:           taskMsg.TaskID,
				TechnicianID: taskMsg.TechnicianID,
				Title:        taskMsg.Title,
				PerformedAt:  time.Now(),
			}
			notification, err := models.NewTaskNotification(task, technician)
			if err != nil {
				log.Printf("Error creating notification: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := c.notificationRepo.Create(ctx, notification); err != nil {
				log.Printf("Error creating notification: %v", err)
				msg.Nack(false, false)
				continue
			}

			msg.Ack(false)
		}
	}()

	return nil
}
