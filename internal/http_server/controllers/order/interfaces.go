package order

type IorderService interface {
	GetOrderById(id string) ([]byte, error)
}
