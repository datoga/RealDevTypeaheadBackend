package main

import (
	"fmt"
	"log"
)

func main() {
	var err error

	config := NewConfig()

	err = config.LoadFromEnv()

	if err != nil {
		panic(err)
	}

	log.Printf("Config: %v\n", config)

	loader := NewDBLoader("names.json")

	popularityModel, err := loader.Load()

	if err != nil {
		panic(err)
	}

	service := NewTypeaheadService(config.SuggestionNumber, popularityModel)

	controller := NewAccountController(config.Port, service)

	fmt.Println("Starting controller at port", config.Port)

	controller.Start()
}
