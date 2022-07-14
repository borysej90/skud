package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type req struct {
	ReaderID int    `json:"reader_id"`
	PassCard string `json:"pass_card"`
}

func main() {
	paths := [][]int{
		{1, 3, 5, 6, 7, 13, 14, 16, 15, 8, 4, 2},
		{8, 9, 10, 1, 9, 11, 16, 12, 10, 3, 4, 2},
		{1, 3, 5, 15, 13, 14, 8, 16, 6, 4, 2, 14},
		{1, 3, 5, 6, 4, 9, 10, 3, 7, 13, 14, 8, 4, 2},
		{1, 9, 11, 12, 10, 3, 7, 16, 6, 6, 15, 8, 4, 2},
	}

	url := "http://" + mustGetEnvVar("APP_HOST") + ":8080/api/access"
	order, err := strconv.Atoi(mustGetEnvVar("ORDER"))
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for i := (order-1)*10 + 1; i <= order*10; i++ {
		i := i
		pathIdx := i % 10
		if pathIdx > 4 {
			pathIdx -= 5
		}
		wg.Add(1)
		go func(empID, pathIdx int) {
			defer wg.Done()
			records := make([][]string, 0)
			defer func() {
				recover()
				if len(records) > 0 {
					saveData(records)
				}
			}()
			for j := 0; j < 10; j++ {
				for _, reader := range paths[pathIdx] {
					reqS := req{
						ReaderID: reader,
						PassCard: fmt.Sprintf("card-%d", empID),
					}
					var reqB []byte
					reqB, err = json.Marshal(reqS)
					if err != nil {
						panic(err)
					}
					start := time.Now().UTC()
					body, err := http.Post(url, "application/json", bytes.NewReader(reqB))
					if err != nil {
						panic(err)
					}
					end := time.Now().UTC()
					success := "True"
					if body.StatusCode != http.StatusOK {
						success = "False"
					}
					records = append(records, []string{strconv.Itoa(order), strconv.Itoa(reader), start.Format(time.RFC3339), end.Format(time.RFC3339), end.Sub(start).String(), success})
				}
			}
		}(i, pathIdx)
	}
	wg.Wait()
}

func saveData(records [][]string) {
	f, err := os.OpenFile("./results/remote-10x20.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	err = w.WriteAll(records)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func mustGetEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(name + " is not set")
	}
	return value
}
