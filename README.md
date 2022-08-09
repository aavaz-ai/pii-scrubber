# PII-Scrubber

## Concepts

### Entity
Entity represents an identifiable piece of text which we are interested in, e.g Date, Email, Credit-Card Number, etc.

### EntityScrubber
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
	Mask(detectedEntity []byte, config EntityConfig) []byte
}
```
**Match** function takes the text as input and returns all the locations for the Entity in the text

**Mask** function is responsible for masking a detected instance of an Entity

## Using Go-Lang Library

To install the library, run the following command in your go project
> go get github.com/aavaz-ai/pii-scrubber

The `Scrubber` interface exposes two high-level functions
- **ScrubTexts**: Useful in scrubbing PII out of the string data
- **ScrubStruct**: An abstraction written on top of ScrubTexts which makes it easier to scrub PII from various specified fields of an object

### Scrub PII from String

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

### Scrub PII from Objects

example:
```go
    type Address struct {
        Location string `pii:"true"`
        ZipCode string `pii:"true"`
    }

	type User struct {
		Name         string            `pii:"true"`
		CustomAttributes           map[string]string `pii:"true"`
		Age           int
		Position           string              `pii:"true"`
		Address           *Address `pii:"true"`
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

	scrubber, _ := piiscrubber.NewWithCustomEntityScrubber(piiscrubber.NewWithCustomEntityScrubberParams{
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
```

## Using CLI
