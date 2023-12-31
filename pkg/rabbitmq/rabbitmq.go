package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wizh/pkg/viper"
	"wizh/pkg/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn   *amqp.Connection
	config = viper.InitConf("rabbitmq")
	Mqurl  = fmt.Sprintf("amqp://%s:%s@%s:%d/%v",
		config.Viper.GetString("server.username"),
		config.Viper.GetString("server.password"),
		config.Viper.GetString("server.host"),
		config.Viper.GetInt("server.port"),
		config.Viper.GetString("server.vhost"))
)

type RabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	QueueName     string // 队列名称
	Exchange      string // 交换机名称
	Key           string // bind key 名称
	Mqurl         string
	Queue         amqp.Queue
	notifyClose   chan *amqp.Error       // 如果异常关闭，会收到数据
	notifyConfirm chan amqp.Confirmation // 消息发送成功确认，会接收到数据
	prefetchCount int
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string, prefetchCount int) *RabbitMQ {
	return &RabbitMQ{
		QueueName:     queueName,
		Exchange:      exchange,
		Key:           key,
		prefetchCount: prefetchCount,
	}
}

// 断开 channel 和 connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string, autoAck bool) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "", config.Viper.GetInt("server.prefetch_count"))

	var err error
	// 获取 connection
	rabbitmq.conn, err = amqp.Dial(Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")

	// 获取 channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel!")

	if !autoAck {
		// 创建一个 qos 控制
		err = rabbitmq.channel.Qos(rabbitmq.prefetchCount, 0, false)
		rabbitmq.failOnErr(err, "failed to create a qos!")
	}

	rabbitmq.channel.NotifyClose(rabbitmq.notifyClose)
	rabbitmq.channel.NotifyPublish(rabbitmq.notifyConfirm)
	return rabbitmq
}

// 生产者
func (r *RabbitMQ) PublishSimple(ctx context.Context, message []byte) error {
	logger := zap.InitLogger()
	// 申请队列，如果队列不存在会自动创建，存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞处理
		nil,   // 额外的属性
	)
	if err != nil {
		logger.Errorln(err)
		if r.conn.IsClosed() || r.channel.IsClosed() {
			return errors.New("RabbitMQ 断开连接，需要重连: " + err.Error())
		}
		return err
	}

	// 调用 channel 发送消息到队列中
	err = r.channel.PublishWithContext(
		ctx,
		r.Exchange,
		r.QueueName,
		false, // 如果为 true，根据 exchange 类型和 route key 规则，如果无法找到符合条件的队列，那么会把发送的消息返回给发送者
		false, // 如果为 true，当 exchange 发送消息到队列后发现队列上没有绑定消费者，则会把消息返回给发送者
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		logger.Errorln(err)
		if r.conn.IsClosed() || r.channel.IsClosed() {
			return errors.New("RabbitMQ 断开连接，需要重连: " + err.Error())
		}
		return err
	}
	return nil
}

// 消费者
func (r *RabbitMQ) ConsumeSimple() (<-chan amqp.Delivery, error) {
	logger := zap.InitLogger()
	// 申请队列
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞处理
		nil,   // 额外的属性
	)
	if err != nil {
		logger.Errorln(err)
	}

	// 接收消息
	msgs, err := r.channel.Consume(
		q.Name, // 队列名称
		"",     // 用来区分多个消费者
		config.Viper.GetBool("consumer.favorite.autoAck"), // 是否自动应答
		false, // 是否独有
		false, // 设置为 true，表示不能将同一个 Connection 中生产者发送的消息传递给这个 Connection 中的消费者
		false, // 列是否阻塞
		nil,
	)

	if err != nil {
		logger.Errorf("RabbitMQ 消费者错误: %v\n", err)
		return nil, err
	}
	return msgs, nil
}
