package main

import (
	db_connection "consumer/db_connection"
	memcache "consumer/memcache"
	nats_streaming_connect "consumer/nats_streaming_connect"
	server "consumer/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 1) Подключение к Postgresql
	db, err := db_connection.DBConnect()
	if err != nil {
		log.Fatalf("Ошибка соединения с базой данных: %v", err)
	}

	// 2) Инициализация кеша cache и заполнение данными из Postgresql

	var cache = memcache.New(5*time.Minute, 10*time.Minute)
	cache.Input(db)

	// 3) Подключение к NATS Streaming серверу
	go nats_streaming_connect.СonnectingNats(db, cache)

	log.Println("Consumer запущен. Ожидание сообщений...")

	// 4) Запуск сервера

	server.ServerStart(cache)

	// Ожидание сигнала для завершения работы приложения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Consumer завершает работу.")

}

// InsertOrder вставляет данные о заказе в базу данных
