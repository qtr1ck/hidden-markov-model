package hmm

import (
	"fmt"
	"github.com/buger/jsonparser"
	"math"
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
			model.AddState(string(key))
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

func (hmms *HiddenMarkovModels) EvaluateModels() {
	greatestProb := 0.0
	bestModel := HiddenMarkovModel{}
	for _, model := range hmms.models {
		modelProb := model.ForwardAlgorithm(&hmms.observations)
		if greatestProb < modelProb {
			greatestProb = modelProb
			bestModel = model
		}
	}
	fmt.Println("Best Model:")
	fmt.Printf("%+v\n", bestModel)
	fmt.Println("Probability:", greatestProb)
	fmt.Println("Log-Probability:", math.Log(greatestProb))

}

func (hmms *HiddenMarkovModels) AddModel(model HiddenMarkovModel) {
	hmms.models = append(hmms.models, model)
}

type HiddenMarkovModel struct {
	states                       []string
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
	// "R->R->0.2" - string format of transition probability
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
	// "R->H->0.2" - string format of emission probability
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

func (hmm *HiddenMarkovModel) ForwardAlgorithm(observations *[]string) float64 {
	// initial step
	alphas := hmm.initForwardAlgo(observations)

	// recursive step
	for t := 1; t < len(*observations); t++ {
		for _, state := range hmm.states {
			sum := 0.0
			for prevState, arr := range alphas {
				sum += arr[t-1] * hmm.stateTransitionProbabilities[prevState][state]
			}
			alphas[state] = append(alphas[state], sum*hmm.emissionProbabilities[state][(*observations)[t]])
		}
	}

	totalProb := 0.0
	// termination step
	for _, alpha := range alphas {
		totalProb += alpha[len(alpha)-1]
	}

	return totalProb
}

func (hmm *HiddenMarkovModel) initForwardAlgo(observations *[]string) map[string][]float64 {
	alpha := make(map[string][]float64)
	for _, state := range hmm.states {
		alpha[state] = append(alpha[state], hmm.initialProbabilities[state]*hmm.emissionProbabilities[state][(*observations)[0]])
	}
	return alpha
}

func (hmm *HiddenMarkovModel) AddState(state string) {
	hmm.states = append(hmm.states, state)
}
