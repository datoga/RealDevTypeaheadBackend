package main

type PopularityModel map[string]int

type Result struct {
	Name  string `json:"name"`
	Times int    `json:"times"`
}
