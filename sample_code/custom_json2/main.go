package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID          string    `json:"id"`
	Items       []Item    `json:"items"`
	DateOrdered time.Time `json:"date_ordered"`
	CustomerID  string    `json:"customer_id"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type Dup Order

	tmp := struct {
		DateOrdered string `json:"date_ordered"`
		Dup
	}{
		Dup: Dup(o),
	}
	tmp.DateOrdered = o.DateOrdered.Format(time.RFC822Z)
	b, err := json.Marshal(tmp)
	return b, err
}

func (o *Order) UnmarshalJSON(b []byte) error {
	type Dup Order

	tmp := struct {
		DateOrdered string `json:"date_ordered"`
		*Dup
	}{
		Dup: (*Dup)(o),
	}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	o.DateOrdered, err = time.Parse(time.RFC822Z, tmp.DateOrdered)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	data := `
	{
		"id": "12345",
		"items": [
			{
				"id": "xyz123",
				"name": "Thing 1"
			},
			{
				"id": "abc789",
				"name": "Thing 2"
			}
		],
		"date_ordered": "01 May 20 13:01 +0000",
		"customer_id": "3"
	}`

	var o Order
	err := json.Unmarshal([]byte(data), &o)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", o)
	fmt.Println(o.DateOrdered.Month())
	out, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
