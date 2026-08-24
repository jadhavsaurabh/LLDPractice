/*
We have to implement a function "PerformDataOperations". Please find below the specifics:
Signature:
func PerformDataOperations(tablesName []string) bool
{
}

Specification/Steps:
Step1:
It should perform read operations on the input tables list.
The data read from the tables is of type string
All the read operations are independant of each other

Step2:
After all the data from the tables is read, we will process the data and prepare "result".
This "result" will be of type string.

	E.g: For 2 tables , table1 and table2. Lets say we get records from these tables as "record1", "record2" respectively

Than the result will be:
result ="record1record2"

Step3: Once the data is processed and we have the "result", we will write this data back to input tables list.
Each write operation is independant of each other.

Step4: Once the write operations are complete, the function will return

	    true: If all the write operations succeeded
	false: If any write operation failed

Note: For the read and write operations, we don't need to write actual db query.
We can write the stubs for these.

Sample stubs:

	func ReadDataFromTable(tableName string) string {
	    return "data from:"+tableName
	}

	func WriteDataToTable(tableName string, data string) error {
	    return nil
	}

Note: These are only samples for reference and are not complete.
You will have to write the function signature as per your logic
*/
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Record struct {
	id   int
	data string
}

func ReadDataFromTable(tableName string) string {
	return "data:" + tableName
}

func WriteDataToTable(tableName string, data string) error {
	if tableName == "tb3" {
		fmt.Println("Failed to write data in table " + tableName)
		return fmt.Errorf("Failed to write data")
	}
	fmt.Println("Writing data to")
	fmt.Println(tableName + "," + data)
	return nil
}

func PerformDataOperations(tablesName []string) bool {
	len := len(tablesName)
	var wg sync.WaitGroup
	dataChan := make(chan Record, len)

	for i, table := range tablesName {
		wg.Add(1)
		go func() {
			defer wg.Done()

			fmt.Printf("\nReading from table %s", table)
			record := ReadDataFromTable(table)
			dataChan <- Record{id: i, data: record}
		}()
	}

	wg.Wait()
	close(dataChan)

	// Combine data
	aggregatedResult := make([]string, len)
	for data := range dataChan {
		aggregatedResult[data.id] = data.data
	}
	res := ""
	for _, aggr := range aggregatedResult {
		res += aggr
	}
	fmt.Println("\nAggregated res: " + res)

	var success atomic.Bool
	success.Store(true)
	for _, table := range tablesName {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := WriteDataToTable(table, res)
			if err != nil {
				success.Store(false)
			}
		}()
	}

	wg.Wait()
	return success.Load()
}

func main() {
	tables := []string{"tb1", "tb2", "tb3", "tb4", "tb5"}
	status := PerformDataOperations(tables)
	fmt.Println("Done with processing: ")
	fmt.Println(status)
}
