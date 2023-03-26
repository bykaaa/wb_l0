package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bykaaa/wb_l0/internal/http_server/models"
)

type orderRepo struct {
	db *sql.DB
}

func (o *orderRepo) GetAllOrders() ([]*models.Order, error) {
	orders := []*models.Order{}
	rows, err := o.db.Query(`
	SELECT o.order_uid , o.track_number ,o.entry ,o.locale ,o.internal_signature ,o.customer_id ,
	o.delivery_service ,o.shardkey ,o.sm_id ,o.date_created ,o.oof_shard ,d.name "delivery_name", d.phone ,d.zip ,
	d.city ,d.address ,d.region ,d.email , p."transaction" ,p.request_id ,p.currency ,p.provider ,p.amount ,p.payment_dt ,p.bank ,
	p.delivery_cost ,p.goods_total ,p.custom_fee ,i.chrt_id ,i.track_number "items_track_number" ,i.price ,i.rid ,i."name" "items_name" ,i.sale ,
	i."size" ,i.total_price ,i.nm_id ,i.brand ,i.status 
	FROM  "Order" o
	LEFT JOIN "Items" i ON i.order_uid =o.order_uid 
	LEFT JOIN "Payment" p ON o.order_uid =p."transaction" 
	LEFT JOIN "Delivery" d ON o.order_uid = d.order_uid order by o.order_uid`)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var prevId string

	for rows.Next() {
		var order models.Order
		var item models.Item
		order.Items = []*models.Item{&item}

		if err := rows.Scan(
			&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId,
			&order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank,
			&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		); err != nil {
			log.Println(err.Error())
		}

		if prevId == *order.OrderUid {
			length := len(orders)
			orders[length-1].Items = append(orders[length-1].Items, &item)
		} else {
			orders = append(orders, &order)
		}

		prevId = *order.OrderUid
		// ! если предыдущий равен текущему id то добавляем items

	}

	return orders, nil
}

func (o *orderRepo) InsertOrders(orders []*models.Order) error {
	a := 0
	var deliveryVals, itemsVals, orderVals, paymentVals []interface{}
	orderStr := `INSERT INTO "Order" (order_uid, track_number, entry, locale,
		internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES`
	deliveryStr := `INSERT INTO "Delivery" (order_uid, name, phone, zip, city, address, region, email) VALUES`
	itemsStr := `INSERT INTO "Items" (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES`

	paymentStr := `INSERT INTO "Payment" (transaction, request_id, currency, provider,
		 amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES`

	quriesMap := map[string]*string{"order": &orderStr, "delivery": &deliveryStr, "items": &itemsStr, "payment": &paymentStr}
	valuesMap := map[string]*[]interface{}{"order": &orderVals, "delivery": &deliveryVals, "items": &itemsVals, "payment": &paymentVals}

	for i, order := range orders {
		orderStr += getValuesFormatPostgreStr(i, 11)
		orderVals = append(orderVals, order.OrderUid, order.TrackNumber, order.Entry,
			order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService,
			order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

		deliveryStr += getValuesFormatPostgreStr(i, 8)
		deliveryVals = append(deliveryVals, order.OrderUid, order.Delivery.Name,
			order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
			order.Delivery.Region, order.Delivery.Email)

		for a < len(order.Items) {
			itemsStr += getValuesFormatPostgreStr(i+a, 12)
			a++
		}

		for _, item := range order.Items {
			itemsVals = append(itemsVals, order.OrderUid, item.ChrtId, item.TrackNumber,
				item.Price, item.Rid, item.Name, item.Sale, item.Size,
				item.TotalPrice, item.NmId, item.Brand, item.Status)
		}

		paymentStr += getValuesFormatPostgreStr(i, 10)
		paymentVals = append(paymentVals, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency,
			order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
			order.Payment.GoodsTotal, order.Payment.CustomFee)

	}

	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, key := range []string{"order", "delivery", "items", "payment"} {
		query := quriesMap[key]
		fmt.Println("KEY:", key, ";", "Query", query)
		q := strings.TrimSuffix(*query, ",")
		values := valuesMap[key]
		if _, err := tx.Exec(q, *values...); err != nil {
			fmt.Println("Insert Error ", key, err)
			return err
		}

	}

	return tx.Commit()
}

func NewOrderRepo(db *sql.DB) *orderRepo {
	return &orderRepo{db: db}
}

func getValuesFormatPostgreStr(i int, length int) string {
	var str = "("
	for j := 0; j < length; j++ {
		str += "$" + strconv.Itoa((i*length)+1+j) + ","
	}

	str = strings.TrimSuffix(str, ",") + "),"

	return str
}
