package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	slurpy "github.com/bobby/slurpy/sdk"
)

func main() {
	// Create a new Slurpy client with logging enabled
	client, err := slurpy.New(slurpy.Config{
		Namespace: "example-app",
		Enabled:   true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Making HTTP requests with Slurpy logging...")

	// Make some example requests
	makeExampleRequests(client)

	fmt.Println("\nRequests logged! Run 'slurpy' CLI to view them.")
}

func makeExampleRequests(client *slurpy.Client) {
	// GET request
	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Printf("GET error: %v\n", err)
	} else {
		fmt.Printf("GET: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	time.Sleep(500 * time.Millisecond)

	// POST request
	resp, err = client.Post("https://httpbin.org/post", "application/json",
		strings.NewReader(`{"message": "Hello from Slurpy!"}`))
	if err != nil {
		fmt.Printf("POST error: %v\n", err)
	} else {
		fmt.Printf("POST: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	time.Sleep(500 * time.Millisecond)

	// PUT request
	resp, err = client.Put("https://httpbin.org/put", "application/json",
		strings.NewReader(`{"updated": true}`))
	if err != nil {
		fmt.Printf("PUT error: %v\n", err)
	} else {
		fmt.Printf("PUT: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	time.Sleep(500 * time.Millisecond)

	// Request that will fail (to show error logging)
	resp, err = client.Get("https://this-domain-does-not-exist-12345.com")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	} else {
		resp.Body.Close()
	}
}
