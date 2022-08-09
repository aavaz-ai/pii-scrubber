package main

import (
	"fmt"
	"strings"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
)

type sampleData struct {
	texts   []string
	pattern string
}

func main() {
	scrubber, _ := piiscrubber.NewDefaultScrubber()

	entityToData := map[piiscrubber.Entity]sampleData{
		piiscrubber.StreetAddress: {
			texts:   addrList,
			pattern: "<STREET_ADDRESS>",
		},
		piiscrubber.CreditCard: {
			texts:   creditCardList,
			pattern: "<CREDIT_CARD>",
		},
		piiscrubber.Email: {
			texts:   emailData,
			pattern: "<EMAIL_ADDRESS>",
		},
	}

	for entity, data := range entityToData {
		response, err := scrubber.ScrubTexts(data.texts)
		if err != nil {
			panic(err)
		}
		total := len(data.texts)
		caught := 0
		for _, scrubbedText := range response {
			if strings.Contains(scrubbedText, data.pattern) {
				caught++
			}
		}

		fmt.Printf("_______%v_______\n", entity)
		fmt.Printf("Total: %v\nCaught: %v\nCoverage: %v%v\n\n=====================\n", total, caught, (caught*100.0)/total, "%")

	}
}
