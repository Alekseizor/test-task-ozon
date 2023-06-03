package url_generation

import "math/rand"

const (
	lenID = 10
)

var (
	dictionaryID = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"
)

func GenerationURL() string {
	var ID string
	for i := 0; i < lenID; i++ {
		ID += string(dictionaryID[rand.Intn(len(dictionaryID))])
	}
	return ID
}
