package db_connection

import (
	"database/sql"
	"fmt"
	"log"
	"encoding/json"
	"os"
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
	configFile, err := os.Open("config.json")
    if err != nil {
        fmt.Println("Error opening config file:", err)
        return nil, err 
    }
    defer configFile.Close()

    // Декодирование данных из файла JSON в структуру Config
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
