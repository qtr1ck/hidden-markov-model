package main

import hmm "hidden-markov-model/src"

func main() {
	inputFile := "resources/models2.json"
	markovModels := hmm.HiddenMarkovModels{}

	err := markovModels.ReadModelsFromFile(inputFile)
	if err != nil {
		panic(err.Error())
	}

	err = markovModels.ReadObservationsFromFile(inputFile)
	if err != nil {
		panic(err.Error())
	}

	markovModels.EvaluateModels()
}

/*
 Das Einlesen der Definition von (beliebig vielen) HMMs soll über Files möglich sein.
 Es soll ebenfalls über Files möglich sein, Sequenzen von Beobachtungen einzulesen.
 Der Forward Algorithmus soll verwendet werden, um zu bestimmen, welches der
eingegebenen HMMs am ehesten zu den gegebenen Beobachtungen passt. Die Wahrscheinlichkeiten,
mit welcher die HMMs und die Beobachtungen „zusammenpassen“, sollen ebenso ausgegeben werden
(sowohl als Wahrscheinlichkeit als auch als log-likelihood).
*/
