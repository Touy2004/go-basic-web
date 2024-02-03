//decode of json data
package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	Id            int
	EmployingName string
	Tel           string
	Address       string
}

func main() {
	data, _ := json.Marshal(&employee{
		Id:            1,
		EmployingName: "John",
		Tel:           "123-4567",
		Address:       "Tokyo",
	})

	fmt.Println(string(data))
}
