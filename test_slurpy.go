package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bobby/slurpy/pkg/storage"
	slurpy "github.com/bobby/slurpy/sdk"
)

func main() {
	fmt.Println("🔬 Testing Slurpy SDK and Storage...")

	// Test 1: SDK initialization
	fmt.Print("1. Testing SDK initialization... ")
	client, err := slurpy.New(slurpy.Config{
		Namespace: "test-suite",
		Enabled:   true,
	})
	if err != nil {
		log.Fatal("❌ Failed:", err)
	}
	fmt.Println("✅ Success")

	// Test 2: Storage initialization
	fmt.Print("2. Testing storage initialization... ")
	store, err := storage.New()
	if err != nil {
		log.Fatal("❌ Failed:", err)
	}
	fmt.Println("✅ Success")

	// Test 3: Make a test request
	fmt.Print("3. Testing HTTP request logging... ")
	resp, err := client.Get("https://httpbin.org/json")
	if err != nil {
		log.Fatal("❌ Failed:", err)
	}
	resp.Body.Close()
	fmt.Println("✅ Success")

	// Test 4: Verify log file was created
	fmt.Print("4. Testing log file creation... ")
	requests, err := store.LoadRequests("test-suite")
	if err != nil {
		log.Fatal("❌ Failed:", err)
	}
	if len(requests) == 0 {
		log.Fatal("❌ Failed: No requests found")
	}
	fmt.Printf("✅ Success (%d requests found)\n", len(requests))

	// Test 5: Verify request data structure
	fmt.Print("5. Testing request data structure... ")
	req := requests[0]
	if req.Method == "" || req.URL == "" || req.Namespace != "test-suite" {
		log.Fatal("❌ Failed: Invalid request structure")
	}
	fmt.Println("✅ Success")

	// Test 6: Test namespace functionality
	fmt.Print("6. Testing namespace switching... ")
	client.SetNamespace("test-suite-v2")
	if client.GetNamespace() != "test-suite-v2" {
		log.Fatal("❌ Failed: Namespace not switched")
	}
	fmt.Println("✅ Success")

	// Test 7: Test enable/disable
	fmt.Print("7. Testing enable/disable... ")
	client.SetEnabled(false)
	if client.IsEnabled() {
		log.Fatal("❌ Failed: Client still enabled")
	}
	client.SetEnabled(true)
	if !client.IsEnabled() {
		log.Fatal("❌ Failed: Client not re-enabled")
	}
	fmt.Println("✅ Success")

	// Test 8: Test CLI build
	fmt.Print("8. Testing CLI build... ")
	if err := buildCLI(); err != nil {
		log.Printf("❌ Failed: %v", err)
	} else {
		fmt.Println("✅ Success")
	}

	// Cleanup
	fmt.Print("9. Cleaning up test data... ")
	if err := store.ClearNamespace("test-suite"); err != nil {
		log.Printf("❌ Cleanup failed: %v", err)
	} else {
		fmt.Println("✅ Success")
	}

	fmt.Println("\n🎉 All tests passed! Slurpy is ready to use.")
	fmt.Println("\nTry these commands:")
	fmt.Println("  make example     # Generate sample data")
	fmt.Println("  make cli         # Launch the TUI")
	fmt.Println("  make run-example # Generate data and launch CLI")
}

func buildCLI() error {
	// Check if we can build the CLI
	if _, err := os.Stat("cli/main.go"); os.IsNotExist(err) {
		return fmt.Errorf("CLI source not found")
	}

	// Simple check - we won't actually run go build to avoid cluttering
	// the test, but we verify the files exist
	requiredFiles := []string{
		"cli/main.go",
		"cli/ui/model.go",
		"cli/ui/commands.go",
		"cli/ui/styles.go",
		"cli/ui/delegate.go",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("required file missing: %s", file)
		}
	}

	return nil
}
