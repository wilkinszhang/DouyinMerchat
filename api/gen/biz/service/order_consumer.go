package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/model"
	"DouyinMerchant/api/gen/conf"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudwego/kitex/pkg/klog"
)

func StartOrderConsumer() {
	dsn := fmt.Sprintf(conf.GetConf().RocketMQ.DSN)
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("ORDER_GROUP"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{dsn})),
	)

	err := c.Subscribe("ORDER_DELAY_TOPIC", consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				orderID := string(msg.Body)
				handleExpiredOrder(orderID)
			}
			return consumer.ConsumeSuccess, nil
		})

	if err != nil {
		klog.Fatalf("消费者启动失败: %v", err)
	}
	c.Start()
}

func handleExpiredOrder(orderID string) {
	tx := mysql.DB.Begin()
	defer tx.Rollback()

	// 获取订单
	order, err := model.GetOrderByID(tx, orderID)
	if err != nil {
		klog.Errorf("查询订单失败: %v", err)
		return
	}

	// 只有placed状态的订单需要关闭
	if order.OrderState == model.OrderStatePlaced {
		err = model.UpdateOrderState(tx, context.Background(), order.UserId, orderID, model.OrderStateCanceled)
		if err != nil {
			klog.Errorf("更新订单状态失败: %v", err)
			return
		}
		klog.Infof("订单已自动关闭: %s", orderID)
	}

	tx.Commit()
}
