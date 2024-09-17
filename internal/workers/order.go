package workers

import (
	"context"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/services"
	"go.uber.org/zap"
	"time"
)

type OrderWorker struct {
	ctx          context.Context
	orderService services.OrderRepository
	orderQueue   chan string
	logger       *zap.Logger
}

func NewOrderWorker(ctx context.Context, orderService services.OrderRepository, orderQueue chan string, logger *zap.Logger) *OrderWorker {
	return &OrderWorker{ctx: ctx, orderService: orderService, orderQueue: orderQueue, logger: logger}
}

func (w *OrderWorker) Run() {
	go w.fillQueue()
	for i := 1; i <= config.Cfg.WorkerCount; i++ {
		go w.worker()
	}
}

func (w *OrderWorker) fillQueue() {
	orders, _ := w.orderService.GetUnComplete(w.ctx)
	for _, orderNum := range orders {
		w.orderQueue <- orderNum
	}
}

func (w *OrderWorker) worker() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case orderNum, ok := <-w.orderQueue:
			if !ok {
				w.logger.Info("order queue is closed")
				return
			}

			result, sleepSeconds, err := client.CheckOrder(orderNum)
			if err != nil {
				w.logger.Info(err.Error())
			} else {
				if err := w.orderService.ChangeStatus(w.ctx, result); err != nil {
					w.logger.Info(err.Error())
				}
			}

			if sleepSeconds > 0 {
				time.Sleep(time.Duration(sleepSeconds) * time.Second)
			}
		}
	}
}
