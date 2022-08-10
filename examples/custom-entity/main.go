package main

import (
	"fmt"
	"regexp"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
)

type orgNameEntityScrubber struct {
}

func (s *orgNameEntityScrubber) Match(text string) [][]int {
	regex := regexp.MustCompile("Enterpret")
	return regex.FindAllStringIndex(text, -1)
}

func (s *orgNameEntityScrubber) Mask(detectedEntity []byte, config *piiscrubber.EntityConfig) []byte {
	return []byte("<ORG_PLACEHOLDER>")
}

func main() {

	texts := []string{
		"Hi this is Anshal, my contact is +919140520809, I am currently working at Enterpret",
	}

	orgNameEntity := piiscrubber.Entity("ORG_NAME")

	scrubber, err := piiscrubber.NewWithCustomEntityScrubbers(piiscrubber.NewWithCustomEntityScrubbersParams{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Phone,
			piiscrubber.Email,
			piiscrubber.SSN,
			orgNameEntity,
		},
		CustomEntityScrubbers: map[piiscrubber.Entity]piiscrubber.EntityScrubber{
			orgNameEntity: &orgNameEntityScrubber{},
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
