package main

import (
	"fmt"
)

func main() {
	fmt.Println("🔍 Testing Basic Java Process Discovery")
	fmt.Println("=====================================")

	// // Create a discoverer
	// discoverer := discovery.NewDiscoverer()

	// // Set up context with timeout
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// // Find Java processes
	// processes, err := discoverer.FindJavaProcesses(ctx)
	// if err != nil {
	// 	log.Fatalf("Error discovering processes: %v", err)
	// }

	// // fmt.Printf("Found %d Java processes:\n\n", len(processes))

	// // Print each process
	// // for i, proc := range processes {
	// // 	fmt.Printf("%d. PID: %d, Owner: %s, Service: %s, JAR: %s\n",
	// // 		i+1, proc.PID, proc.Owner, proc.ServiceName, proc.JarFile)
	// // }

	// if len(processes) == 0 {
	// 	fmt.Println("No Java processes found!")
	// 	fmt.Println("Make sure you have Java applications running.")
	// }
}
