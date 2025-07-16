package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	slurpy "github.com/bobby/slurpy/sdk"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func main() {
	// Example 1: Multiple namespaces
	fmt.Println("=== Advanced Slurpy Example ===")

	// Client for user service
	userClient, err := slurpy.New(slurpy.Config{
		Namespace: "user-service",
		Enabled:   true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Client for API gateway (different namespace)
	gatewayClient, err := slurpy.New(slurpy.Config{
		Namespace: "api-gateway",
		Enabled:   true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Example 2: CRUD operations
	fmt.Println("\n1. Creating user...")
	createUser(userClient)

	time.Sleep(200 * time.Millisecond)

	fmt.Println("2. Reading users...")
	getUsers(userClient)

	time.Sleep(200 * time.Millisecond)

	fmt.Println("3. Updating user...")
	updateUser(userClient)

	time.Sleep(200 * time.Millisecond)

	// Example 3: Different service
	fmt.Println("4. Gateway health check...")
	healthCheck(gatewayClient)

	time.Sleep(200 * time.Millisecond)

	// Example 4: Custom request with headers
	fmt.Println("5. Custom request with headers...")
	customRequest(gatewayClient)

	time.Sleep(200 * time.Millisecond)

	// Example 5: Runtime configuration changes
	fmt.Println("6. Disabling logging temporarily...")
	userClient.SetEnabled(false)
	silentRequest(userClient) // This won't be logged

	fmt.Println("7. Re-enabling logging...")
	userClient.SetEnabled(true)
	getUsers(userClient) // This will be logged

	// Example 6: Namespace switching
	fmt.Println("8. Switching namespace...")
	userClient.SetNamespace("user-service-v2")
	getUsers(userClient) // Logged under new namespace

	fmt.Println("\n=== Example complete! ===")
	fmt.Println("Check the CLI to see requests organized by namespace:")
	fmt.Println("- user-service")
	fmt.Println("- api-gateway")
	fmt.Println("- user-service-v2")
}

func createUser(client *slurpy.Client) {
	user := User{Name: "Alice", Role: "admin"}
	body, _ := json.Marshal(user)

	resp, err := client.Post("https://httpbin.org/post", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Create error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Created user: %d\n", resp.StatusCode)
}

func getUsers(client *slurpy.Client) {
	resp, err := client.Get("https://httpbin.org/get?service=users")
	if err != nil {
		fmt.Printf("Get error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Retrieved users: %d\n", resp.StatusCode)
}

func updateUser(client *slurpy.Client) {
	user := User{ID: 123, Name: "Alice Updated", Role: "super-admin"}
	body, _ := json.Marshal(user)

	resp, err := client.Put("https://httpbin.org/put", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Update error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Updated user: %d\n", resp.StatusCode)
}

func healthCheck(client *slurpy.Client) {
	resp, err := client.Get("https://httpbin.org/status/200")
	if err != nil {
		fmt.Printf("Health check error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Health check: %d\n", resp.StatusCode)
}

func customRequest(client *slurpy.Client) {
	req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	if err != nil {
		fmt.Printf("Request creation error: %v\n", err)
		return
	}

	// Add custom headers
	req.Header.Set("X-API-Key", "secret-key-123")
	req.Header.Set("X-Request-ID", "req-456")
	req.Header.Set("User-Agent", "Slurpy-Advanced-Example/1.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Custom request error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Custom headers: %d\n", resp.StatusCode)
}

func silentRequest(client *slurpy.Client) {
	resp, err := client.Get("https://httpbin.org/get?silent=true")
	if err != nil {
		fmt.Printf("Silent error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Silent request (not logged): %d\n", resp.StatusCode)
}
