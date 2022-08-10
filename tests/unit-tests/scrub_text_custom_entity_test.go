package test

import (
	"regexp"
	"testing"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
	"github.com/stretchr/testify/assert"
)

type customTestEntityScrubber struct {
}

func (s *customTestEntityScrubber) Match(text string) [][]int {
	regex := regexp.MustCompile("Aavaz")
	return regex.FindAllStringIndex(text, -1)
}

func (s *customTestEntityScrubber) Mask(detectedEntity []byte, config *piiscrubber.EntityConfig) []byte {
	return []byte("Enterpret")
}

func Test_ScrubTextsCustom(t *testing.T) {
	texts := []string{
		"Hi this is Anshal with, +919140520809, Working at Aavaz",
	}

	expectedTexts := []string{
		"Hi this is Anshal with, <PHONE_NUMBER>, Working at Enterpret",
	}

	scrubber, err := piiscrubber.NewWithCustomEntityScrubbers(piiscrubber.NewWithCustomEntityScrubbersParams{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Phone,
			piiscrubber.Email,
			piiscrubber.SSN,
			"COMPANY_NAME",
		},
		Config: map[piiscrubber.Entity]*piiscrubber.EntityConfig{
			piiscrubber.Email: {
				MaskWithChar:         runePtr('x'),
				UnmaskedSuffixOffset: 4,
			},
			piiscrubber.CreditCard: {
				ReplaceWith: stringPtr("<CREDIT_CARD_DETECTED>"),
			},
			"COMPANY_NAME": &piiscrubber.EntityConfig{},
		},
		CustomEntityScrubbers: map[piiscrubber.Entity]piiscrubber.EntityScrubber{
			"COMPANY_NAME": &customTestEntityScrubber{},
		},
	})

	assert.NoError(t, err)

	response, err := scrubber.ScrubTexts(texts)
	assert.NoError(t, err)

	assert.Equal(t, expectedTexts, response)
}
