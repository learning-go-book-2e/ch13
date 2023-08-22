package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	err := ProcessPerson()
	if err != nil {
		slog.Error("error in processPerson", "msg", err)
	}
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ProcessPerson() error {
	toFile := Person{
		Name: "Fred",
		Age:  40,
	}

	// Write it out
	tmpFile, err := os.CreateTemp(os.TempDir(), "sample-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	err = json.NewEncoder(tmpFile).Encode(toFile)
	if err != nil {
		return err
	}
	err = tmpFile.Close()
	if err != nil {
		return err
	}

	// Read it back in again
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	var fromFile Person
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	if err != nil {
		return err
	}
	err = tmpFile2.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", fromFile)
	return nil
}
