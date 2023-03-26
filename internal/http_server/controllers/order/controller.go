package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type getOrderController struct {
	orderService IorderService
}

func (c *getOrderController) GetOrder(f *fiber.Ctx) error {
	query := GetOrderQuery{}
	err := f.QueryParser(&query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	order, err := c.orderService.GetOrderById(query.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.Response().Header.Add("Content-Type", "application/json; charset=utf-8")

	return f.Send(order)
}

func NewOrderController(orderService IorderService) *getOrderController {
	return &getOrderController{orderService: orderService}
}
