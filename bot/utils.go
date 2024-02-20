package bot

import "math/rand"

func shuffleStrings(inputSlice []string) []string {
	rand.Shuffle(len(inputSlice), func(i, j int) {
		inputSlice[i], inputSlice[j] = inputSlice[j], inputSlice[i]
	})
	return inputSlice
}
