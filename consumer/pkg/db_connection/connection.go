package db_connection

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func DBConnect() (*sql.DB, error) {
	// Получаем текущую директорию
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return nil, err
	}

	// Открываем файл config.json относительно текущей директории
	configFile, err := os.Open(filepath.Join(dir, "/pkg/db_connection/config.json"))
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return nil, err
	}
	defer configFile.Close()

	// Остальной код остается без изменений
	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return nil, err
	}
	var connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Успешное подключение к базе данных PostgreSQL!")
	return db, nil
}
