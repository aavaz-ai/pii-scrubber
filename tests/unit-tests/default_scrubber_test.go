package test

import (
	"fmt"
	"testing"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
	"github.com/stretchr/testify/assert"
)

func Test_DefaultScrubber(t *testing.T) {
	texts := []string{
		"Hi this is Anshal with, +919140520809",
		"Hi ping mein at anshaldwivedi@gmail.com",
		"here, 6011553157232994",
		"My SSN is488-23-3729. Details can be found at https://aavaz.ai/emp/488-23-3729",
		"My email is 9144520109@gmail.com",
	}

	expectedTexts := []string{
		"Hi this is Anshal with, <PHONE_NUMBER>",
		"Hi ping mein at <EMAIL_ADDRESS>",
		"here, <CREDIT_CARD>",
		"My SSN is<US_SSN>. Details can be found at https://aavaz.ai/emp/488-23-3729",
		"My email is <EMAIL_ADDRESS>",
	}

	scrubber, _ := piiscrubber.NewDefaultScrubber()

	response, err := scrubber.ScrubTexts(texts)
	fmt.Println(response)
	assert.NoError(t, err)

	assert.Equal(t, expectedTexts, response)
}
