package main

import (
	"L0/server"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"os"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go server.CreateAndLaunchApi()
	go server.LaunchStreaming()
	wg.Wait()
}

func publish() {
	conn, _ := stan.Connect("test-cluster", "str-1", stan.NatsURL("0.0.0.0:4223"))
	defer conn.Close()
	file, _ := os.Open("./model.json")
	data, _ := ioutil.ReadAll(file)
	conn.Publish("order-channel", data)
	println("successfully publish")
}
