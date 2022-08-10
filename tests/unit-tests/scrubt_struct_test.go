package test

import (
	"testing"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
	"github.com/stretchr/testify/assert"
)

func Test_ScrubStruct(t *testing.T) {

	type nestedSampleStruct struct {
		K string
	}

	type sampleStruct struct {
		Val         string            `pii:"true"`
		M           map[string]string `pii:"true"`
		I           int
		S           string              `pii:"true"`
		N           *nestedSampleStruct `pii:"true"`
		EmailPII    string              `pii:"true"`
		EmailNonPII string
	}

	v := sampleStruct{
		Val: "Anshal +9140528009",
		M: map[string]string{
			"PIIKey": "Hello here is my credit card 6011553157232994",
		},
		I: 10,
		S: "Hello!",
		N: &nestedSampleStruct{
			K: "My 488-23-3729",
		},
		EmailPII:    "abc@gmail.com",
		EmailNonPII: "abc@enterpret.com",
	}

	scrubber, _ := piiscrubber.NewWithCustomEntityScrubbers(piiscrubber.NewWithCustomEntityScrubbersParams{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Phone,
			piiscrubber.Email,
			piiscrubber.SSN,
			"COMPANY_NAME",
		},
		Config: map[piiscrubber.Entity]*piiscrubber.EntityConfig{
			piiscrubber.CreditCard: {
				ReplaceWith: stringPtr("CREDIT_CARD_DETECTED"),
			},
			"COMPANY_NAME": {},
		},
		CustomEntityScrubbers: map[piiscrubber.Entity]piiscrubber.EntityScrubber{
			"COMPANY_NAME": &customTestEntityScrubber{},
		},
	})

	expectedResponse := sampleStruct{
		Val: "Anshal <PHONE_NUMBER>",
		M: map[string]string{
			"PIIKey": "Hello here is my credit card CREDIT_CARD_DETECTED",
		},
		I: 10,
		S: "Hello!",
		N: &nestedSampleStruct{
			K: "My <US_SSN>",
		},
		EmailPII:    "<EMAIL_ADDRESS>",
		EmailNonPII: "abc@enterpret.com",
	}

	response, err := scrubber.ScrubStruct(v)
	assert.NoError(t, err)

	assert.Equal(t, expectedResponse, response)
}

type customTestEntityScrubberError struct {
}

func (s *customTestEntityScrubberError) Match(text string) [][]int {
	return [][]int{{0}}
}

func (s *customTestEntityScrubberError) Mask(detectedEntity []byte, config *piiscrubber.EntityConfig) []byte {
	return []byte("Enterpret")
}

func Test_ScrubStruct_Failure(t *testing.T) {

	type nestedSampleStruct struct {
		K string
	}

	type sampleStruct struct {
		Val         string            `pii:"true"`
		M           map[string]string `pii:"true"`
		I           int
		S           string              `pii:"true"`
		N           *nestedSampleStruct `pii:"true"`
		EmailPII    string              `pii:"true"`
		EmailNonPII string
	}

	v := sampleStruct{
		Val: "Anshal +9140520009",
		M: map[string]string{
			"PIIKey": "Hello here is my credit card 6011553157232994",
		},
		I: 10,
		S: "Hello!",
		N: &nestedSampleStruct{
			K: "My 488-23-3729",
		},
		EmailPII:    "abc@gmail.com",
		EmailNonPII: "abc@enterpret.com",
	}

	scrubber, _ := piiscrubber.NewWithCustomEntityScrubbers(piiscrubber.NewWithCustomEntityScrubbersParams{
		BlacklistedEntities: []piiscrubber.Entity{
			piiscrubber.CreditCard,
			piiscrubber.Phone,
			piiscrubber.Email,
			piiscrubber.SSN,
			"COMPANY_NAME",
		},
		Config: map[piiscrubber.Entity]*piiscrubber.EntityConfig{
			piiscrubber.CreditCard: {
				ReplaceWith: stringPtr("CREDIT_CARD_DETECTED"),
			},
			"COMPANY_NAME": &piiscrubber.EntityConfig{},
		},
		CustomEntityScrubbers: map[piiscrubber.Entity]piiscrubber.EntityScrubber{
			"COMPANY_NAME": &customTestEntityScrubberError{},
		},
	})

	response, err := scrubber.ScrubStruct(v)
	assert.Error(t, err)
	assert.Nil(t, response)
}
