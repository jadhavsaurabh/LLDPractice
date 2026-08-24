package main

import (
	"fmt"
	"sync"
)

type ConnectionManager struct {
	conn int
}

var (
	dbCon *ConnectionManager
	once  sync.Once
)

func GetConnectionManager() *ConnectionManager {
	once.Do(func() {
		fmt.Println("establishing db connection")
		dbCon = &ConnectionManager{
			conn: 1,
		}
	})

	return dbCon
}

func main() {
	var wg sync.WaitGroup

	for num := range 20 {
		wg.Add(1)

		go func(num int) {
			defer wg.Done()

			db := GetConnectionManager()

			fmt.Println("Calling query", num, "connection:", db.conn)
		}(num)
	}

	wg.Wait()
}
