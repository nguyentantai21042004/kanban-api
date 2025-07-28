package consumer

import (
	smtpConsumer "gitlab.com/tantai-kanban/kanban-api/internal/core/smtp/rabbitmq/consumer"
	smtpUC "gitlab.com/tantai-kanban/kanban-api/internal/core/smtp/usecase"
)

func (srv Consumer) mapHandlers() error {
	smtpUC := smtpUC.New(srv.l, srv.smtpConfig)

	var forever chan bool
	smtpConsumer.NewConsumer(srv.l, srv.amqpConn, smtpUC).Consume()
	<-forever

	return nil
}
