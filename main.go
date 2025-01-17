package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"hrbac/global"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

type Person struct {
	Name     string  `nebula:"name"`
	Age      int     `nebula:"age"`
	Likeness float64 `nebula:"likeness"`
}

func main() {
	prepareSpace()
	NewSessionPool()
	defer global.SessionPool.Close()

	// execute query
	{
		insertVertexes := "INSERT VERTEX person(name, age) VALUES " +
			"'Bob':('Bob', 10), " +
			"'Lily':('Lily', 9), " +
			"'Tom':('Tom', 10), " +
			"'Jerry':('Jerry', 13), " +
			"'John':('John', 11);"

		// Insert multiple vertexes
		resultSet, err := global.SessionPool.Execute(insertVertexes)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		CheckResultSet(insertVertexes, resultSet)
	}
	{
		// Insert multiple edges
		insertEdges := "INSERT EDGE like(likeness) VALUES " +
			"'Bob'->'Lily':(80.0), " +
			"'Bob'->'Tom':(70.0), " +
			"'Lily'->'Jerry':(84.0), " +
			"'Tom'->'Jerry':(68.3), " +
			"'Bob'->'John':(97.2);"

		resultSet, err := global.SessionPool.Execute(insertEdges)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		CheckResultSet(insertEdges, resultSet)
	}
	// Extract data from the resultSet
	var err error
	{
		query := "GO FROM 'Bob' OVER like YIELD $^.person.name AS name, $^.person.age AS age, like.likeness AS likeness"
		// Send query in goroutine
		wg := sync.WaitGroup{}
		wg.Add(1)
		var resultSet *nebula.ResultSet
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			resultSet, err = global.SessionPool.Execute(query)
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			CheckResultSet(query, resultSet)
			var personList []Person
			resultSet.Scan(&personList)
			fmt.Printf("personList: %v\n", personList)
			// personList: [{Bob 10 97.2} {Bob 10 80} {Bob 10 70}]
		}(&wg)
		wg.Wait()

		// Get all column names from the resultSet
		colNames := resultSet.GetColNames()
		fmt.Printf("column names: %s\n", strings.Join(colNames, ", "))

		// Get a row from resultSet
		record, err := resultSet.GetRowValuesByIndex(0)
		if err != nil {
			log.Println(err.Error())
		}
		// Print whole row
		fmt.Printf("row elements: %s\n", record.String())
		// Get a value in the row by column index
		valueWrapper, err := record.GetValueByIndex(0)
		if err != nil {
			log.Println(err.Error())
		}
		// Get type of the value
		fmt.Printf("valueWrapper type: %s \n", valueWrapper.GetType())
		// Check if valueWrapper is a string type
		if valueWrapper.IsString() {
			// Convert valueWrapper to a string value
			v1Str, err := valueWrapper.AsString()
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Printf("Result of ValueWrapper.AsString(): %s\n", v1Str)
		}
		// Print ValueWrapper using String()
		fmt.Printf("Print using ValueWrapper.String(): %s", valueWrapper.String())
	}
	// Drop space
	{
		query := fmt.Sprintf(`DROP SPACE IF EXISTS %s`, global.SPACE)
		// Send query
		resultSet, err := global.SessionPool.Execute(query)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		CheckResultSet(query, resultSet)
	}
	fmt.Print("\n")
	log.Println("Nebula Go Client Session Pool Example Finished")
}
