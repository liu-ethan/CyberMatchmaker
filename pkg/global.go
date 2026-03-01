/**
 * @author 刘潇翰
 * @since 2026/3/1
 */
package pkg

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
	MQ    *amqp.Connection
)
