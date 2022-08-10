# Concepts

## Entity
Entity represents an identifiable piece of text which we are interested in, e.g Date, Email, Credit-Card Number, etc.

## EntityScrubber
EntityScrubber responsible for detecting and masking an entity in the provided input. The library provides pre-built scrubbers for the following entities

	Date          Entity = "DATE"
	Time          Entity = "TIME"
	CreditCard    Entity = "CREDIT_CARD"
	Phone         Entity = "PHONE"
	Link          Entity = "LINK"
	Email         Entity = "EMAIL"
	IP            Entity = "IP"
	NotKnownPort  Entity = "UNKNOWN_PORT"
	BtcAddress    Entity = "BTC_ADDRESS"
	StreetAddress Entity = "STREET_ADDRESS"
	ZipCode       Entity = "ZIP_CODE"
	PoBox         Entity = "PO_BOX"
	SSN           Entity = "SSN"
	MD5Hex        Entity = "MD5_HEX"
	SHA1Hex       Entity = "SHA1_HEX"
	SHA256Hex     Entity = "SHA_256_HEX"
	GUID          Entity = "GUID"
	ISBN          Entity = "ISBN"
	MACAddress    Entity = "MAC_ADDRESS"
	IBAN          Entity = "IBAN"
	GitRepo       Entity = "GIT_REPO"
	StrictLink    Entity = "STRICT_LINK"


User can override the implementation of any of the existing scrubber or add their own custom-entities and corresponding scrubbers. 
Entity scrubber implements the following interface

```go
type EntityScrubber interface {
	Match(text string) [][]int
	Mask(detectedEntity []byte, config *EntityConfig) []byte
}
```
**Match** function takes the text as input and returns all the locations for the Entity in the text

**Mask** function is responsible for masking a detected instance of an Entity

# Installation
To install the library, run the following command in your go project
> go get github.com/aavaz-ai/pii-scrubber
<br></br>

# Usage
The `Scrubber` interface exposes two high-level functions
- **ScrubTexts**: Useful in scrubbing PII out of the string data
- **ScrubStruct**: An abstraction written on top of ScrubTexts which makes it easier to scrub PII from various specified fields of an object

## Scrub PII from String

example:
```go

	texts := []string{
		"Hi my phone number is +919140520809",
	}

	scrubber, err := piiscrubber.NewDefaultScrubber()
    if err != nil {
        panic(err)
    }

	response, err := scrubber.ScrubTexts(texts)
    if err != nil {
        panic(err)
    }

    fmt.Println(response)
```
Output:
```go
["Hi my phone number is <PHONE_NUMBER>"]
```

## Scrub PII from Objects

example:
```go

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
```
Output:
```json
{
  "Name": "Anshal <PHONE_NUMBER>",
  "CustomAttributes": {
    "PIIKey": "Hello here is my credit card <CREDIT_CARD>"
  },
  "Age": 10,
  "Position": "Software Engineer",
  "Address": {
    "Location": "My <US_SSN>",
    "ZipCode": "<ZIP_CODE>"
  },
  "Email": "<EMAIL_ADDRESS>"
}
```

# Advance Usage

## [ Add a Custom Entity ](https://github.com/aavaz-ai/pii-scrubber/tree/master/examples/custom-entity)

In the following example, we implement an orgNameEntityScrubber, that matches a certain organisation's name, and masks it with a placeholder value

``` go
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
```
Output:
```json
["Hi this is Anshal, my contact is <PHONE_NUMBER>, I am currently working at <ORG_PLACEHOLDER>"]
```
<br></br>
## [ Override an Entity-Scrubber ](https://github.com/aavaz-ai/pii-scrubber/tree/master/examples/override-entity-scrubber)
```go
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
```
Output:
```json
["Hi this is Anshal, my credit card is 4263982640269299, and <CUSTOM_CREDIT_CARD>, I am currently working at Enterpret"]
```
<br></br>
## [ Use Config for Custom Masking ](https://github.com/aavaz-ai/pii-scrubber/tree/master/examples/mask-with-config)
EntityConfig provides limited parameters to customise the masking operation for a detected entity. More advance requirements can be addressed by overriding the EntityScrubber itself

```go
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
```
Output:
```json
["Hi this is Anshal, my contact is <PHONE_NUMBER>, and credit card is XXXXXXXXXXXX9299, I am currently working at Enterpret"]
```

***

<br></br>

# Tests

## Unit Tests
``` bash
cd tests/unit-tests
go test ./...
```

## Coverage Tests
```bash
cd tests/benchmarks/coverage
go run ./...
```
Output:
```text
_______STREET_ADDRESS_______
Total: 200
Caught: 162
Coverage: 81%

=====================
_______CREDIT_CARD_______
Total: 26
Caught: 23
Coverage: 88%

=====================
_______EMAIL_______
Total: 13
Caught: 13
Coverage: 100%

=====================
```

## Performance Tests
```bash
cd tests/benchmarks/performance
go test -bench=.
```
Output:
```bash
goos: darwin
goarch: amd64
pkg: github.com/aavaz-ai/pii-scrubber/tests/benchmarks/performance
cpu: VirtualApple @ 2.50GHz
Benchmark_1000Sentences-10            87          12112771 ns/op
Benchmark_100Sentences-10            543           2194837 ns/op
```