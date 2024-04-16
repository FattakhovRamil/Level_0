package nats_streaming_connect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	// "strconv"
	"time"

	memcache "consumer/memcache"
	. "consumer/order_struct"

	"github.com/nats-io/nats.go"
)

func СonnectingNats(db *sql.DB, c *memcache.Cache) error {
	fmt.Println("Установка соединения с сервером NATS Streaming")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer nc.Close()
	var order Order
	subscription, err := nc.Subscribe("test.subject", func(msg *nats.Msg) {
		fmt.Println("Установка соединения с сервером NATS Streaming - успешно")
		log.Printf("Получено сообщение: %s", string(msg.Data))

		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatalf("Ошибка парсинга JSON: %v", err)

		}

		// Отправка заказа в базу данных
		err := InsertOrder(db, order, c)
		if err != nil {
			log.Fatalf("Ошибка записи заказа в базу данных: %v", err)
			return
		}

	})

	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()
	select {}
	//go printCachePeriodically(c, time.Second*15)
}

func InsertOrder(db *sql.DB, order Order, c *memcache.Cache) error {

	fmt.Println("вставляет данные о заказе в базу данных")

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Откатить транзакцию при любой ошибке

	// Вставка данных в таблицу orders
	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSign, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OOFShard)
	if err != nil {
		return err
	}

	// Получаем идентификатор заказа
	var orderID int

	err = tx.QueryRow("SELECT MAX(id) FROM orders WHERE order_uid = $1", order.OrderUID).Scan(&orderID)
	if err != nil {
		return err
	}

	// Вставка данных в таблицу delivery
	_, err = tx.Exec(`
		INSERT INTO delivery (order_id, name, phone, zip, city, address, region, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	// Вставка данных в таблицу payment
	_, err = tx.Exec(`
		INSERT INTO payment (order_id, transaction, request_id, currency, provider, 
		amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	// Вставка данных в таблицу items
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (order_id, chrt_id, track_number, price, rid, name, sale, size, 
			total_price, nm_id, brand, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func printCachePeriodically(cache *memcache.Cache, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Вызов метода печати кэша
			cache.PrintCache()
		}
	}
}

func PrintMap(m map[string]interface{}) {
	fmt.Println("Printing Map:")
	for key, value := range m {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
}
