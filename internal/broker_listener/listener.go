package broker_listener

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/bykaaa/wb_l0/internal/http_server/models"
	"github.com/bykaaa/wb_l0/pkg/validator"
)

type orderRepo interface {
	InsertOrders(orders []*models.Order) error
}

type brokerListener struct {
	orderRepo   orderRepo
	ordersCache *map[string][]byte
	pendingOrders
}

type pendingOrders struct {
	mu               sync.Mutex
	notRecordOrders  []*models.Order
	pendingMsgs      []*stan.Msg
	pendingOrdersMap map[string]struct{}
}

func (bl *brokerListener) Run() {

	opts := []stan.Option{stan.NatsURL("nats://nats-streaming-server:4222")}
	sc, err := stan.Connect("test-cluster", "2", opts...)
	if err != nil {
		log.Fatal(err)
	}

	msgChan := make(chan *stan.Msg)
	ordersSub, err := sc.Subscribe("orders", func(m *stan.Msg) {
		msgChan <- m
	}, stan.DurableName("order_service"), stan.SetManualAckMode(), stan.MaxInflight(1000))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("run")

	go bl.listen(ordersSub, msgChan)
	go bl.defferedInsertOrders()
}

func NewBrokerListener(ordersCache *map[string][]byte, ororderRepo orderRepo) *brokerListener {
	return &brokerListener{ordersCache: ordersCache, orderRepo: ororderRepo, pendingOrders: pendingOrders{
		pendingOrdersMap: make(map[string]struct{}, 1000),
	}}

}

func (bl *brokerListener) listen(ordersSub stan.Subscription, msgChan <-chan *stan.Msg) {

	fmt.Println("listen")
	defer ordersSub.Unsubscribe()

	for msg := range msgChan {
		order, err := validator.ValidatBodyAndGetData[*models.Order](msg.Data)
		fmt.Println("New msg", order.Payment.RequestId)
		if err != nil {
			fmt.Println("Error", err)
			msg.Ack()
			continue
		}

		bl.pendingOrders.mu.Lock()

		if *order.OrderUid != *order.Payment.Transaction {
			bl.pendingOrders.mu.Unlock()
			msg.Ack()
			continue
		}

		if _, ok := (*bl.ordersCache)[*order.OrderUid]; ok {
			bl.pendingOrders.mu.Unlock()
			msg.Ack()
			continue
		}

		if _, ok := bl.pendingOrdersMap[*order.OrderUid]; ok {
			bl.pendingOrders.mu.Unlock()
			msg.Ack()
			continue
		}

		bl.pendingOrdersMap[*order.OrderUid] = struct{}{}
		bl.pendingOrders.notRecordOrders = append(bl.pendingOrders.notRecordOrders, order)
		bl.pendingOrders.pendingMsgs = append(bl.pendingOrders.pendingMsgs, msg)
		bl.pendingOrders.mu.Unlock()

	}

}

func (bl *brokerListener) defferedInsertOrders() {
	for {
		time.Sleep(3 * time.Second)
		bl.pendingOrders.mu.Lock()
		// fmt.Println("deffered insert", bl.pendingOrders.notRecordOrders)
		notInsertedOrders := bl.pendingOrders.notRecordOrders
		notAcceptedMsgs := bl.pendingOrders.pendingMsgs
		bl.pendingOrders.notRecordOrders = []*models.Order{}
		bl.pendingOrders.pendingMsgs = []*stan.Msg{}
		bl.pendingOrders.mu.Unlock()

		if len(notInsertedOrders) != 0 {
			err := bl.orderRepo.InsertOrders(notInsertedOrders)
			if err != nil {
				log.Fatal(err)
			}
			bl.pendingOrders.mu.Lock()
			for i, msg := range notAcceptedMsgs {
				bytes, _ := json.Marshal(notInsertedOrders[i])

				(*bl.ordersCache)[*notInsertedOrders[i].OrderUid] = bytes
				msg.Ack()
			}
			bl.pendingOrders.pendingOrdersMap = make(map[string]struct{}, 1000)
			bl.pendingOrders.mu.Unlock()
			fmt.Println("Saved !!!")
		}
	}
}
