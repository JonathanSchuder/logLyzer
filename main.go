package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var lines = make(chan string, 100)

// var mu sync.Mutex

func main() {

	infos, err := os.Stat("./access.log")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(infos.Size())

	file, err := os.Open("./access.log")

	if err != nil {
		fmt.Println(err)
		return
	} else {
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup

	go func() {

		for scanner.Scan() {
			lines <- scanner.Text()
		}

		close(lines)
	}()

	workerCount := 5
	maps := make([]map[int]int, workerCount)

	for i := range workerCount {
		wg.Add(1)
		currMap := make(map[int]int, 50)
		maps[i] = currMap
		go logWorker(&wg, currMap)
	}

	wg.Wait()

	results := make(map[int]int, 1000)

	for _, currMap := range maps {
		for key, value := range currMap {
			results[key] += value
		}
	}

	fmt.Println(results)
}

func logWorker(wg *sync.WaitGroup, workerMap map[int]int) {

	defer wg.Done()

	for line := range lines {
		parts := strings.Fields(line)
		code, err := strconv.Atoi(parts[len(parts)-1])

		if err != nil {
			continue
		}

		workerMap[code]++
	}
}
