package server

import (
	"L0/database"
	"L0/entity"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"sync"
)

var ordersCache = make(map[string]entity.Order, 0)
var ordersCacheByte = make(map[string][]byte, 0)
var ordersIds = make(map[string]string, 0)
var db *sqlx.DB

type Server struct {
	ordersCache map[string]entity.Order
	orderIds    map[string]string
	db          *sqlx.DB
}

func CreateAndLaunchApi() {
	cfg := database.DefaultConfig()
	var err error
	db, err = database.NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}
	loadCacheFromDB(db)
	println(len(ordersCache))

	router := routing.New()
	router.To("GET", "/order/<id>", func(context *routing.Context) error {
		context.Write(getOrderById(context))
		//fmt.Fprintf(context, "ID: %v", context.Param("id"))
		return nil
	})
	fasthttp.ListenAndServe(":8080", router.HandleRequest)

	//router := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	//router := gin.New()
	//router.GET("/order/:id", getOrderByID)
	//router.Run("localhost:8080")
}

func loadCacheFromDB(db *sqlx.DB) {

	orders, err := database.GetAllOrders(db)
	if err != nil {
		log.Fatal(err)
		return
	}
	ordersCache = orders
	for s, order := range ordersCache {
		ordersIds[s] = "+"
		temp, _ := json.Marshal(order)
		ordersCacheByte[order.OrderUID] = temp
	}
	println(len(orders))
}

func LaunchStreaming() {
	conn, err := stan.Connect("test-cluster", "sub-1", stan.NatsURL("0.0.0.0:4223"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	conn.Subscribe("order-channel", ReceiveJson)
	wait()
}

func wait() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func ReceiveJson(message *stan.Msg) {
	var order entity.Order
	json.Unmarshal(message.Data, &order)
	insertToDB(order)
	insertToCache(order)
}

func getOrderById(context *routing.Context) []byte {
	id := context.Param("id")
	if ordersIds[id] != "" {
		context.SetStatusCode(200)
		return ordersCacheByte[id]
	} else {
		context.SetStatusCode(404)
		return []byte{}
	}
}

func getOrderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if ordersIds[id] != "" {
		ctx.JSON(http.StatusOK, ordersCache[id])
	} else {
		getOrderFromDB(id)
	}
}

func getOrderFromDB(id string) (entity.Order, error) {
	return entity.Order{}, nil
}

func insertToDB(order entity.Order) {
	println(db)
	data, _ := json.Marshal(order)
	database.InsertOrder(db, order.OrderUID, data)
	println("insert to db from stream ", order.OrderUID)
}
func insertToCache(order entity.Order) {
	ordersCache[order.OrderUID] = order
	ordersIds[order.OrderUID] = "s"
	temp, _ := json.Marshal(order)
	ordersCacheByte[order.OrderUID] = temp
}
