package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bykaaa/wb_l0/internal/broker_listener"
	"github.com/bykaaa/wb_l0/internal/http_server"
	"github.com/bykaaa/wb_l0/internal/http_server/models"
	"github.com/bykaaa/wb_l0/internal/http_server/repositories"

	_ "github.com/lib/pq"
)

func main() {
	userName, password, dbName, host, port := "postgres", "1234", "wb_l0", "db", 5432
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", userName, password, dbName, host, port)
	log.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	orderRepo := repositories.NewOrderRepo(db)

	orders, err := orderRepo.GetAllOrders()

	ordersCache := initCache(orders)

	s := http_server.InitHttpServer(db, &ordersCache)

	l := broker_listener.NewBrokerListener(&ordersCache, orderRepo)

	l.Run()

	s.Server.Listen(":8072")
}

func initCache(orders []*models.Order) map[string][]byte {
	cache := map[string][]byte{}
	for _, order := range orders {
		bytes, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}

		cache[*order.OrderUid] = bytes
	}
	return cache
}
