//encode to json format
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type employees struct {
	Id            int
	EmployingName string
	Tel           string
	Address       string
}

func main() {
	e := employees{}
	err := json.Unmarshal([]byte(`{"Id": 1, "EmployingName": "John", "Tel": "123-4567", "Address": "Tokyo"}`), &e)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.Address)
}
