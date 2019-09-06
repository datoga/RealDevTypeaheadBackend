package main

import (
	"errors"
	"log"
	"sort"
	"strings"
)

type TypeaheadService struct {
	suggestionNumber int
	popularityModel  PopularityModel
}

func NewTypeaheadService(suggestionNumber int, popularityModel PopularityModel) *TypeaheadService {
	return &TypeaheadService{
		suggestionNumber: suggestionNumber,
		popularityModel:  popularityModel,
	}
}

func (service TypeaheadService) Get(query string) (*Result, bool) {

	log.Println("Querying", query)

	for name, votes := range service.popularityModel {
		if strings.EqualFold(name, query) {
			result := Result{
				Name:  name,
				Times: votes,
			}
			return &result, true
		}
	}

	return nil, false
}

func (service TypeaheadService) model2results(model PopularityModel) []Result {
	results := []Result{}

	for k, v := range model {
		result := Result{
			Name:  k,
			Times: v,
		}

		results = append(results, result)
	}

	return results
}

func (service TypeaheadService) Query(prefix string) []Result {
	hits := PopularityModel{}

	upperPrefix := strings.ToUpper(prefix)

	var upperName string

	for name, votes := range service.popularityModel {
		upperName = strings.ToUpper(name)

		if strings.HasPrefix(upperName, upperPrefix) {
			hits[name] = votes
		}
	}

	results := service.model2results(hits)

	sort.Slice(results, func(i, j int) bool {
		if results[i].Times != results[j].Times {
			return results[i].Times > results[j].Times
		}

		return results[i].Name < results[j].Name
	})

	if len(results) > service.suggestionNumber {
		results = results[:service.suggestionNumber]
	}

	return results
}

func (service *TypeaheadService) Vote(name string) error {
	hit, ok := service.Get(name)

	if !ok {
		return errors.New("Error, " + name + " is not found")
	}

	service.popularityModel[hit.Name]++

	return nil
}
