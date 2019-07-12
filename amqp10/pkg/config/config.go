package config

import "pack.ag/amqp"

type Config struct {
	Session *amqp.Session
	Client  *amqp.Client
}