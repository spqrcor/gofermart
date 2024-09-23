package workers

import (
	"context"
	"fmt"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/services"
	"go.uber.org/zap"
	"time"
)

type OrderWorker struct {
	ctx          context.Context
	orderService *services.OrderService
	orderQueue   chan string
	logger       *zap.Logger
	conf         config.Config
	orderClient  *client.OrderClient
}

func NewOrderWorker(opts ...func(*OrderWorker)) *OrderWorker {
	orderWorker := &OrderWorker{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt(orderWorker)
	}
	return orderWorker
}

func WithCtx(ctx context.Context) func(*OrderWorker) {
	return func(o *OrderWorker) {
		o.ctx = ctx
	}
}

func WithOrderService(orderService *services.OrderService) func(*OrderWorker) {
	return func(o *OrderWorker) {
		o.orderService = orderService
	}
}

func WithOrderQueue(orderQueue chan string) func(*OrderWorker) {
	return func(o *OrderWorker) {
		o.orderQueue = orderQueue
	}
}

func WithLogger(logger *zap.Logger) func(*OrderWorker) {
	return func(o *OrderWorker) {
		o.logger = logger
	}
}

func WithConfig(conf config.Config) func(*OrderWorker) {
	return func(o *OrderWorker) {
		o.conf = conf
	}
}

func WithOrderClient(orderClient *client.OrderClient) func(*OrderWorker) {
	return func(o *OrderWorker) {
		o.orderClient = orderClient
	}
}

func (w *OrderWorker) Run() {
	retryCount := 0
	go w.fillQueue()
	for i := 1; i <= w.conf.WorkerCount; i++ {
		errorCh := w.worker()

		go func() {
			for err := range errorCh {
				if err != nil && retryCount < w.conf.RetryStartWorkerCount {
					w.worker()
					retryCount++
				}
			}
		}()
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
