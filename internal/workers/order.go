package workers

import (
	"context"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/services"
	"time"
)

type OrderWorker struct {
	ctx          context.Context
	orderService services.OrderRepository
}

func NewOrderWorker(ctx context.Context, orderService services.OrderRepository) *OrderWorker {
	return &OrderWorker{ctx: ctx, orderService: orderService}
}

func (w *OrderWorker) Send(orderNum string) error {
	if err := client.SendOrder(orderNum); err != nil {
		return err
	}
	data := services.OrderFromAccrual{Order: orderNum, Status: "PROCESSING"}
	if err := w.orderService.ChangeStatus(w.ctx, data); err != nil {
		return err
	}
	return nil
}

func (w *OrderWorker) Run() {
	//_ = client.SetReward()
	go w.doInterval()
}

func (w *OrderWorker) doInterval() {
	for range time.Tick(time.Second * 5) {
		logger.Log.Info("execute 5s")

		orders, _ := w.orderService.GetUnComplete(w.ctx)
		for _, order := range orders {
			logger.Log.Info("before check")
			result, err := client.CheckOrder(order)
			if err != nil {
				logger.Log.Info(err.Error())

			} else {
				err = w.orderService.ChangeStatus(w.ctx, result)
				if err != nil {
					logger.Log.Info(err.Error())

				}
				logger.Log.Info("after check")
			}
		}
	}
}
