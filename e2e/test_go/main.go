package main

import (
	context "context"
	"fmt"
	"math/rand"
	"os"
	sync "sync"

	"google.golang.org/grpc"
)

// WaitGroup is used to wait for the program to finish goroutines.
var wg sync.WaitGroup

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Sending requests..")

	wg.Add(5)
	go sendRequests(conn, 100)
	go sendRequests(conn, 100)
	go sendRequests(conn, 100)
	go sendRequests(conn, 100)
	go sendRequests(conn, 1000)

	wg.Wait()

	fmt.Println("Program finished")
}

func sendRequests(conn *grpc.ClientConn, numberOfRequests int) {
	client := NewEntrypointClient(conn)

	fails := 0

	for i := 0; i < numberOfRequests; i++ {
		generatedName := fmt.Sprintf("Alex-%d", (rand.Intn(10000)))
		resp, err := client.Greet(context.Background(), &Request{Name: generatedName})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		generatedResponse := fmt.Sprintf("Hello %s!, how are you? from nodeC", generatedName)
		if generatedResponse != resp.Greeting {
			fmt.Println(generatedResponse, "---", resp.Greeting)
			fails++
		}
	}
	fmt.Println("Fails", fails)
	wg.Done()
}
