package main

import (
	"fmt"
	"regexp"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
)

type creditCardOverrideScrubber struct {
}

func (s *creditCardOverrideScrubber) Match(text string) [][]int {
	// implement the logic to detect credit-card number here

	regex := regexp.MustCompile("4263 9826 4026 9299")
	return regex.FindAllStringIndex(text, -1)
}

func (s *creditCardOverrideScrubber) Mask(detectedEntity []byte, config *piiscrubber.EntityConfig) []byte {
	return []byte("<CUSTOM_CREDIT_CARD>")
}

func main() {

	texts := []string{
		"Hi this is Anshal, my credit card is 4263982640269299, and 4263 9826 4026 9299, I am currently working at Enterpret",
	}

	scrubber, err := piiscrubber.NewWithCustomEntityScrubbers(piiscrubber.NewWithCustomEntityScrubbersParams{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Email,
			piiscrubber.SSN,
		},
		CustomEntityScrubbers: map[piiscrubber.Entity]piiscrubber.EntityScrubber{
			piiscrubber.CreditCard: &creditCardOverrideScrubber{},
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
