package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/model"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/order"
	"DouyinMerchant/mq"
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
} // NewPlaceOrderService new PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run create note info
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	// Finish your business logic.
	if len(req.OrderItems) == 0 {
		err = fmt.Errorf("OrderItems empty")
		return
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		orderId, _ := uuid.NewUUID()

		o := &model.Order{
			OrderId:      orderId.String(),
			OrderState:   model.OrderStatePlaced,
			UserId:       req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: model.Consignee{
				Email: req.Email,
			},
		}
		if req.Address != nil {
			a := req.Address
			o.Consignee.Country = a.Country
			o.Consignee.State = a.State
			o.Consignee.City = a.City
			o.Consignee.StreetAddress = a.StreetAddress
		}
		if err := tx.Create(o).Error; err != nil {
			return err
		}

		var itemList []*model.OrderItem
		for _, v := range req.OrderItems {
			itemList = append(itemList, &model.OrderItem{
				OrderIdRefer: o.OrderId,
				ProductId:    v.Item.ProductId,
				Quantity:     v.Item.Quantity,
				Cost:         v.Cost,
			})
		}
		if err := tx.Create(&itemList).Error; err != nil {
			return err
		}
		resp = &order.PlaceOrderResp{
			Order: &order.OrderResult{
				OrderId: orderId.String(),
			},
		}

		msg := &primitive.Message{
			Topic: "ORDER_DELAY_TOPIC",
			Body:  []byte(o.OrderId),
		}
		// 设置30分钟延迟（RocketMQ支持18个等级，这里需要换算）
		msg.WithDelayTimeLevel(16) // 对应30分钟

		_, err = mq.OrderProducer.SendSync(context.Background(), msg)
		if err != nil {
			klog.Errorf("发送延迟消息失败: %v", err)
			// 这里可以添加补偿逻辑（如记录到数据库定时任务）
		}

		return nil
	})

	return
}
