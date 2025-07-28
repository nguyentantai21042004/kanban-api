package rabbitmq

import rabbitmqPkg "gitlab.com/tantai-kanban/kanban-api/pkg/rabbitmq"

const (
	SEND_EMAIL_EXCHANGE_NAME = "smtp_send_email_exc"
	SEND_EMAIL_QUEUE_NAME    = "smtp_send_email"
)

var (
	SendEmailExc = rabbitmqPkg.ExchangeArgs{
		Name:       SEND_EMAIL_EXCHANGE_NAME,
		Type:       rabbitmqPkg.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
)
