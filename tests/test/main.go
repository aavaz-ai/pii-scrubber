package main

import (
	"encoding/json"
	"fmt"

	piiscrubber "github.com/aavaz-ai/pii-scrubber"
)

func main() {
	type Address struct {
		Location string
		ZipCode  string
	}

	type User struct {
		Name             string            `pii:"true"`
		CustomAttributes map[string]string `pii:"true"`
		Age              int
		Position         string
		Address          *Address `pii:"true"`
		Email            string   `pii:"true"`
	}

	v := User{
		Name: "Anshal +9140528009",
		CustomAttributes: map[string]string{
			"PIIKey": "Hello here is my credit card 6011553157232994",
		},
		Age:      10,
		Position: "Software Engineer",
		Address: &Address{
			Location: "My 488-23-3729",
			ZipCode:  "22132",
		},
		Email: "abc@gmail.com",
	}

	scrubber, err := piiscrubber.NewDefaultScrubber()
	if err != nil {
		panic(err)
	}

	out, err := scrubber.ScrubStruct(v)
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(out)
	fmt.Println(string(data))
}
