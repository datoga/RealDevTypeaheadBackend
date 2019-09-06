package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DBLoader struct {
	source string
}

func NewDBLoader(source string) *DBLoader {
	return &DBLoader{
		source: source,
	}
}

func (loader DBLoader) Load() (PopularityModel, error) {
	jsonFile, err := os.Open(loader.source)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var model PopularityModel

	err = json.Unmarshal(byteValue, &model)

	if err != nil {
		return nil, err
	}

	return model, nil
}
