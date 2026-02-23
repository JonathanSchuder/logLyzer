package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	file, err := os.Create("access.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	statusCodes := []int{200, 200, 200, 200, 301, 404, 401, 500}
	endpoints := []string{"/home", "/login", "/api/v1/users", "/assets/logo.png", "/contact"}

	for i := 0; i < 5000; i++ {
		code := statusCodes[rand.Intn(len(statusCodes))]
		path := endpoints[rand.Intn(len(endpoints))]
		timestamp := time.Now().Format("02/Jan/2006:15:04:05")

		line := fmt.Sprintf("127.0.0.1 - - [%s] \"GET %s HTTP/1.1\" %d\n", timestamp, path, code)

		file.WriteString(line)
	}

	fmt.Println("Successfully generated access.log with 5000 lines.")
}
