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
var totalFails int = 0
var numberOfClients int = 5
var numberOfRequestsPerClient = 100 // try different loads!

func increaseTotalRequestCounter() {
	counterMutex.Lock()
	totalRequests++
	counterMutex.Unlock()
}

func increaseTotalFailCounter() {
	counterMutex.Lock()
	totalFails++
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
	fmt.Printf("Total number of fails: %d\n", totalFails)
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

	fmt.Printf("Created client with address: %v+\n", client)

	fails := 0
	myCounter := 0

	for i := 0; i < numberOfRequests; i++ {
		generatedRequest, expectedResponse := generateRequest()

		resp, err := client.Greet(context.Background(), &proto.Request{Name: generatedRequest})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if expectedResponse != resp.Greeting {
			fmt.Println(expectedResponse, "---", resp.Greeting)
			fails++
			increaseTotalFailCounter()
		}

		increaseTotalRequestCounter()
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

func generateRequest() (string, string) {
	var (
		generatedRequest = ""
		expectedResponse = ""
	)

	if randomInt := rand.Intn(20); randomInt == 0 {
		generatedRequest = "early exit"
		expectedResponse = "early exit"
	} else if randomInt == 1 {
		generatedRequest = "early reply"
		expectedResponse = "early reply"
	} else {
		generatedRequest = fmt.Sprintf("Alex-%d", (rand.Intn(10000)))
		expectedResponse = fmt.Sprintf("Hello %s! greetings from nodeA, nodeB and nodeC!", generatedRequest)
	}

	return generatedRequest, expectedResponse
}
