package workers

import (
	"context"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/services"
	"time"
)

const workerCount = 3

type OrderWorker struct {
	ctx          context.Context
	orderService services.OrderRepository
	orderQueue   chan string
}

func NewOrderWorker(ctx context.Context, orderService services.OrderRepository, orderQueue chan string) *OrderWorker {
	return &OrderWorker{ctx: ctx, orderService: orderService, orderQueue: orderQueue}
}

func (w *OrderWorker) Run() {
	go w.fillQueue()
	for i := 1; i <= workerCount; i++ {
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
				logger.Log.Info("order queue is closed")
				return
			}

			result, sleepSeconds, err := client.CheckOrder(orderNum)
			if err != nil {
				logger.Log.Info(err.Error())
			} else {
				err = w.orderService.ChangeStatus(w.ctx, result)
				if err != nil {
					logger.Log.Info(err.Error())
				}
			}

			if sleepSeconds > 0 {
				time.Sleep(time.Duration(sleepSeconds) * time.Second)
			}
		}
	}
}
