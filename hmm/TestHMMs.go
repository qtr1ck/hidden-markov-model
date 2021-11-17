package hmm

import "fmt"

func TestHMMs(filename string) error {
	markovModels := HiddenMarkovModels{}
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Evaluation of the best model of the file", filename)

	err := markovModels.ReadModelsFromFile(filename)
	if err != nil {
		return err
	}

	err = markovModels.ReadObservationsFromFile(filename)
	if err != nil {
		return err
	}

	markovModels.EvaluateModels()
	fmt.Println("--------------------------------------------------------------------------------")
	return err
}
