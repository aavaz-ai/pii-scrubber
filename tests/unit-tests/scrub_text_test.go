package test

import (
	"testing"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
	"github.com/stretchr/testify/assert"
)

func Test_ScrubTexts(t *testing.T) {
	texts := []string{
		"Hi this is Anshal with, +919140520809",
		"Hi ping mein at anshaldwivedi@gmail.com",
		"here, 6011553157232994",
		"My SSN is488-23-3729. Details can be found at https://aavaz.ai/emp/488-23-3729",
		"My email is 9144520109@gmail.com",
	}

	expectedTexts := []string{
		"Hi this is Anshal with, <PHONE_NUMBER>",
		"Hi ping mein at xxxxxxxxxxxxxxxxxxx.com",
		"here, <CREDIT_CARD_DETECTED>",
		"My SSN is<US_SSN>. Details can be found at https://aavaz.ai/emp/488-23-3729",
		"My email is xxxxxxxxxxxxxxxx.com",
	}

	scrubber, _ := piiscrubber.New(piiscrubber.Params{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Phone,
			piiscrubber.Email,
			piiscrubber.SSN,
		},
		IgnoredEntities: []piiscrubber.Entity{
			piiscrubber.StrictLink,
		},
		Config: map[piiscrubber.Entity]*piiscrubber.EntityConfig{
			piiscrubber.Email: {
				MaskWithChar:         runePtr('x'),
				UnmaskedSuffixOffset: 4,
			},
			piiscrubber.CreditCard: {
				ReplaceWith: stringPtr("<CREDIT_CARD_DETECTED>"),
			},
		},
	})

	response, err := scrubber.ScrubTexts(texts)
	assert.NoError(t, err)

	assert.Equal(t, expectedTexts, response)
}

func runePtr(r rune) *rune {
	return &r
}

func stringPtr(s string) *string {
	return &s
}
