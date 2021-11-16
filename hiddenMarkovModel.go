package main

import (
	"github.com/buger/jsonparser"
	"os"
	"strconv"
	"strings"
)

type HiddenMarkovModels struct {
	observations []string
	models       []HiddenMarkovModel
}

func (hmms *HiddenMarkovModels) ReadModelsFromFile(filename string) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		model := HiddenMarkovModel{}
		// parse initial probabilities
		err = jsonparser.ObjectEach(value, func(key []byte, innerValue []byte, innerDataType jsonparser.ValueType, innerOffset int) error {
			parsedProb, err := strconv.ParseFloat(string(innerValue), 64)
			model.AddInitProb(string(key), parsedProb)
			return err
		}, "initialProbabilities")

		if err != nil {
			return
		}

		// parse transition probabilities
		transProb, _ := jsonparser.GetString(value, "stateTransitionProbabilities")
		err = model.AddTransProbFromJSONString(transProb)
		if err != nil {
			return
		}
		// parse emission probabilities
		emProb, _ := jsonparser.GetString(value, "emissionProbabilities")
		err = model.AddEmProbFromJSONString(emProb)
		if err != nil {
			return
		}
		hmms.AddModel(model)
	}, "models")

	return err
}

func (hmms *HiddenMarkovModels) ReadObservationsFromFile(filename string) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	obsStr, err := jsonparser.GetString(data, "observations")
	if err != nil {
		return err
	}

	observations := strings.Split(obsStr, ",")
	for _, obs := range observations {
		hmms.observations = append(hmms.observations, strings.TrimSpace(obs))
	}
	return nil
}

func (hmms *HiddenMarkovModels) AddModel(model HiddenMarkovModel) {
	hmms.models = append(hmms.models, model)
}

type HiddenMarkovModel struct {
	// states                       []string  TODO: necessary?
	initialProbabilities         map[string]float64
	stateTransitionProbabilities map[string]map[string]float64
	emissionProbabilities        map[string]map[string]float64
}

func (hmm *HiddenMarkovModel) AddInitProb(key string, probability float64) {
	if hmm.initialProbabilities == nil {
		hmm.initialProbabilities = make(map[string]float64)
	}
	hmm.initialProbabilities[key] = probability
}

func (hmm *HiddenMarkovModel) AddTransProbFromJSONString(probString string) error {
	// "R->R->0.2, R->S->0.1, R->C->0.7, S->R->0.3, S->S->0.4, S->C->0.3, C->R->0.1, C->S->0.4, C->C->0.5"
	probs := strings.Split(probString, ",")
	for _, prob := range probs {
		sp := strings.Split(prob, "->")
		parsedProb, err := strconv.ParseFloat(strings.TrimSpace(sp[2]), 64)
		hmm.AddTransProb(sp[0], sp[1], parsedProb)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hmm *HiddenMarkovModel) AddEmProbFromJSONString(probString string) error {
	probs := strings.Split(probString, ",")
	for _, prob := range probs {
		sp := strings.Split(prob, "->")
		parsedProb, err := strconv.ParseFloat(strings.TrimSpace(sp[2]), 64)
		hmm.AddEmProb(sp[0], sp[1], parsedProb)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hmm *HiddenMarkovModel) AddTransProb(from string, to string, prob float64) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if hmm.stateTransitionProbabilities == nil {
		hmm.stateTransitionProbabilities = make(map[string]map[string]float64)
	}
	if hmm.stateTransitionProbabilities[from] == nil {
		hmm.stateTransitionProbabilities[from] = make(map[string]float64)
	}
	hmm.stateTransitionProbabilities[from][to] = prob
}

func (hmm *HiddenMarkovModel) AddEmProb(from string, to string, prob float64) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if hmm.emissionProbabilities == nil {
		hmm.emissionProbabilities = make(map[string]map[string]float64)
	}
	if hmm.emissionProbabilities[from] == nil {
		hmm.emissionProbabilities[from] = make(map[string]float64)
	}
	hmm.emissionProbabilities[from][to] = prob
}