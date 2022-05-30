package main

import (
	context "context"
	"fmt"
	"main/proto"
	"math/rand"
	"os"
	sync "sync"

	"google.golang.org/grpc"
)

// WaitGroup is used to wait for the program to finish goroutines.
var wg sync.WaitGroup

func main() {

	fmt.Println("Sending requests..")

	wg.Add(4)
	go sendRequests(100)
	go sendRequests(100)
	go sendRequests(100)
	go sendRequests(100)
	// go sendRequests(conn, 100)
	// go sendRequests(conn, 100)
	// go sendRequests(conn, 1000)

	wg.Wait()

	fmt.Println("Program finished")
}

func sendRequests(numberOfRequests int) {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	client := proto.NewEntrypointClient(conn)
	myExpectedResponses := make(map[string]bool)

	fmt.Printf("%v+\n", client)

	fails := 0

	for i := 0; i < numberOfRequests; i++ {
		generatedName := fmt.Sprintf("Alex-%d", (rand.Intn(10000)))
		generatedResponse := fmt.Sprintf("Hello %s!, how are you? from nodeC", generatedName)
		myExpectedResponses[generatedResponse] = true

		resp, err := client.Greet(context.Background(), &proto.Request{Name: generatedName})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if generatedResponse != resp.Greeting && !myExpectedResponses[resp.Greeting] {
			fmt.Println(generatedResponse, "---", resp.Greeting)
			fails++
		}
	}
	fmt.Println("Fails", fails)
	wg.Done()
}
