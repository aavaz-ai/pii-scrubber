package main

import (
	"fmt"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
)

func main() {

	texts := []string{
		"Hi this is Anshal, my contact is +919140520809, and credit card is 4263982640269299, I am currently working at Enterpret",
	}

	scrubber, err := piiscrubber.NewWithCustomEntityScrubbers(piiscrubber.NewWithCustomEntityScrubbersParams{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Phone,
			piiscrubber.Email,
			piiscrubber.SSN,
		},
		Config: map[piiscrubber.Entity]*piiscrubber.EntityConfig{
			piiscrubber.CreditCard: {
				UnmaskedSuffixOffset: 4,
				MaskWithChar:         runePtr('X'),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	response, err := scrubber.ScrubTexts(texts)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}

func runePtr(r rune) *rune {
	return &r
}
