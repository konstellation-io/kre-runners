package main

import (
	context "context"
	"fmt"
	"main/proto"
	"math"
	"math/rand"
	"os"
	sync "sync"
	"time"

	"google.golang.org/grpc"
)

// WaitGroup is used to wait for the program to finish goroutines.
var wg sync.WaitGroup
var counterMutex sync.Mutex
var totalRequests int = 0
var numberOfClients int = 5
var numberOfRequestsPerClient = 100

func increaseCounter() {
	counterMutex.Lock()
	totalRequests++
	counterMutex.Unlock()
}

func main() {

	fmt.Println("Sending requests..")

	wg.Add(numberOfClients)

	for i := 0; i < numberOfClients; i++ {
		go sendRequests(numberOfRequestsPerClient)
	}

	wg.Wait()

	fmt.Printf("Total number of requests: %d\n", totalRequests)
	fmt.Println("Program finished")
}

func sendRequests(numberOfRequests int) {
	start := time.Now()
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	client := proto.NewEntrypointClient(conn)
	myExpectedResponses := make(map[string]bool)

	fmt.Printf("%v+\n", client)

	fails := 0
	myCounter := 0

	for i := 0; i < numberOfRequests; i++ {

		generatedName := fmt.Sprintf("Alex-%d", (rand.Intn(10000)))
		generatedResponse := fmt.Sprintf("Hello %s!, how are you? from nodeC", generatedName)
		// fmt.Printf("Send message to %s", generatedName)
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

		increaseCounter()
		myCounter++
		myPercentage := float64(myCounter) / float64(numberOfRequests) * 100
		if math.Mod(myPercentage, 10) == 0 {
			fmt.Printf("%v%%\n", myPercentage)
		}
	}

	elapsed := time.Since(start)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Sent %d messages with execution time %s with %d fails\n", numberOfRequests, elapsed, fails)
	wg.Done()
}
