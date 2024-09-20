package workers

import (
	"context"
	"fmt"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/services"
	"go.uber.org/zap"
	"time"
)

type OrderWorker struct {
	ctx          context.Context
	orderService *services.OrderService
	orderQueue   chan string
	logger       *zap.Logger
	workerCount  int
	orderClient  *client.OrderClient
}

func NewOrderWorker(ctx context.Context, orderService *services.OrderService, orderQueue chan string, logger *zap.Logger, workerCount int, orderClient *client.OrderClient) *OrderWorker {
	return &OrderWorker{ctx: ctx, orderService: orderService, orderQueue: orderQueue, logger: logger, workerCount: workerCount, orderClient: orderClient}
}

func (w *OrderWorker) Run() {
	go w.fillQueue()
	for i := 1; i <= w.workerCount; i++ {
		w.worker()
	}
}

func (w *OrderWorker) fillQueue() {
	orders, _ := w.orderService.GetUnComplete(w.ctx)
	for _, orderNum := range orders {
		w.orderQueue <- orderNum
	}
}

func (w *OrderWorker) worker() chan error {
	errors := make(chan error)
	go func() {
		defer close(errors)
		for {
			select {
			case <-w.ctx.Done():
				return
			case orderNum, ok := <-w.orderQueue:
				if !ok {
					w.logger.Info("order queue is closed")
					errors <- fmt.Errorf("order queue is closed")
					return
				}

				result, sleepSeconds, err := w.orderClient.CheckOrder(orderNum)
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
	}()
	return errors
}
