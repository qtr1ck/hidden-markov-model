package main

import (
	"hidden-markov-model/hmm"
	"path"
)

func main() {
	testfiles := [3]string{"models1.json", "models2.json", "models3.json"}
	for _, file := range testfiles {
		err := hmm.TestHMMs(path.Join("resources", file))
		if err != nil {
			panic(err.Error())
		}
	}
}
