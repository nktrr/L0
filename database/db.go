package database

import (
	"L0/entity"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
	db, err := sqlx.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DefaultConfig() Config {
	cfg := Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "4eJ7sKwB",
		DBName:   "lab",
		SSLMode:  "disable",
	}
	return cfg
}

func GetAllOrders(db *sqlx.DB) (map[string]entity.Order, error) {
	var data [][]byte
	query := `SELECT "order" FROM orders`
	if err := db.Select(&data, query); err != nil {
		return nil, err
	}
	orders := make(map[string]entity.Order)
	for _, datum := range data {
		println("select ", len(datum))
		temp := entity.Order{}
		json.Unmarshal(datum, &temp)
		orders[temp.OrderUID] = temp
	}
	return orders, nil

	//orders := make(map[string]entity.Order)
	//for _, datum := range data {
	//	var temp entity.Order
	//	json.Unmarshal(datum.Order, &temp)
	//	orders[temp.OrderUID] = temp
	//	//json.Unmarshal(datum.Order, &orders)
	//}
	return orders, nil
}

func InsertOrder(db *sqlx.DB, id string, data []byte) {
	println("insert ", len(data))
	//query := fmt.Sprintf(`INSERT INTO orders("id","order") VALUES ('%v', '%v')`, id, data)
	//query = strings.Replace(query, "[", "{", -1)
	//query = strings.Replace(query, "]", "}", -1)
	_, err := db.Exec("INSERT INTO orders (\"id\", \"order\") VALUES ($1,$2::bytea)", id, data)
	if err != nil {
		log.Fatal(err)
	}
}

type DBOrder struct {
	Id    string
	Order []byte
}
