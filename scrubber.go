package piiscrubber

import (
	"fmt"
	"sort"

	goworker "github.com/anshal21/go-worker"
)

// Scrubber ...
type Scrubber interface {
	ScrubTexts(texts []string) ([]string, error)
	ScrubStruct(obj interface{}) (interface{}, error)
}

// Params ...
type Params struct {
	BlacklistedEntities []Entity
	IgnoredEntities     []Entity
	Config              map[Entity]*EntityConfig
}

// New DefaultScrubber ...
func NewDefaultScrubber() (Scrubber, error) {
	return &scrubber{
		blacklistedEntities: []Entity{
			StreetAddress,
			CreditCard,
			Phone,
			Email,
			IP,
			ZipCode,
			PoBox,
			SSN,
			ISBN,
			MACAddress,
			IBAN,
		},
		ignoredEntities: []Entity{
			StrictLink,
			GitRepo,
		},
	}, nil
}

// NewScrubber ...
func New(params Params) (Scrubber, error) {
	for key, val := range params.Config {
		if err := val.isValid(); err != nil {
			return nil, fmt.Errorf("in config for entity: %v, error: %v", key, err.Error())
		}
	}

	return &scrubber{
		blacklistedEntities: params.BlacklistedEntities,
		ignoredEntities:     params.IgnoredEntities,
		config:              params.Config,
	}, nil
}

// NewWithCustomEntityScrubberParams ...
type NewWithCustomEntityScrubberParams struct {
	BlacklistedEntities   []Entity
	IgnoredEntities       []Entity
	Config                map[Entity]*EntityConfig
	CustomEntityScrubbers map[Entity]EntityScrubber
}

var (
	// ErrInvalidMatchIndices ...
	ErrInvalidMatchIndices = fmt.Errorf("invalid match generated by the entity scrubber")
)

// NewWithCustomEntityScrubber ...
func NewWithCustomEntityScrubber(params NewWithCustomEntityScrubberParams) (Scrubber, error) {
	for key := range params.CustomEntityScrubbers {
		if _, ok := params.Config[key]; !ok {
			return nil, fmt.Errorf("missing config for custom or overriden entity: %v", key)
		}
	}

	for key, val := range params.Config {
		// skip the check for custom scrubbers
		if _, ok := params.CustomEntityScrubbers[key]; ok {
			continue
		}
		if err := val.isValid(); err != nil {
			return nil, fmt.Errorf("in config for entity: %v, error: %v", key, err.Error())
		}
	}

	return &scrubber{
		blacklistedEntities:   params.BlacklistedEntities,
		ignoredEntities:       params.IgnoredEntities,
		config:                params.Config,
		userProvidedScrubbers: params.CustomEntityScrubbers,
	}, nil
}

type scrubber struct {
	config                map[Entity]*EntityConfig
	blacklistedEntities   []Entity
	ignoredEntities       []Entity
	userProvidedScrubbers map[Entity]EntityScrubber
}

// Entity ...
type Entity string

// Possible Entities ...
const (
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
)

func getPlaceholderValue(piiPlaceholder string) *string {
	str := fmt.Sprintf("<%v>", piiPlaceholder)
	return &str
}

var (
	_defaultEntityConfigs = map[Entity]*EntityConfig{
		Date:          {ReplaceWith: getPlaceholderValue("DATE")},
		Time:          {ReplaceWith: getPlaceholderValue("TIME")},
		Phone:         {ReplaceWith: getPlaceholderValue("PHONE_NUMBER")},
		Link:          {ReplaceWith: getPlaceholderValue("LINK")},
		Email:         {ReplaceWith: getPlaceholderValue("EMAIL_ADDRESS")},
		IP:            {ReplaceWith: getPlaceholderValue("IP")},
		NotKnownPort:  {ReplaceWith: getPlaceholderValue("NOT_KNOWN_PORT")},
		CreditCard:    {ReplaceWith: getPlaceholderValue("CREDIT_CARD")},
		BtcAddress:    {ReplaceWith: getPlaceholderValue("BITCOIN_ADDRESS")},
		StreetAddress: {ReplaceWith: getPlaceholderValue("STREET_ADDRESS")},
		ZipCode:       {ReplaceWith: getPlaceholderValue("ZIP_CODE")},
		PoBox:         {ReplaceWith: getPlaceholderValue("PO_BOX")},
		SSN:           {ReplaceWith: getPlaceholderValue("US_SSN")},
		MD5Hex:        {ReplaceWith: getPlaceholderValue("MD5_HEX")},
		SHA1Hex:       {ReplaceWith: getPlaceholderValue("SHA1_HEX")},
		SHA256Hex:     {ReplaceWith: getPlaceholderValue("SHA_256_HEX")},
		GUID:          {ReplaceWith: getPlaceholderValue("GUID")},
		ISBN:          {ReplaceWith: getPlaceholderValue("ISBN")},
		MACAddress:    {ReplaceWith: getPlaceholderValue("MAC_ADDRESS")},
		IBAN:          {ReplaceWith: getPlaceholderValue("IBAN")},
		GitRepo:       {ReplaceWith: getPlaceholderValue("GIT_REPO")},
		StrictLink:    {ReplaceWith: getPlaceholderValue("STRICT_LINK")},
	}

	_defaultEntityScrubbers = map[Entity]EntityScrubber{
		MACAddress:    &mACAddressEntityScrubber{},
		MD5Hex:        &mD5HexEntityScrubber{},
		ISBN:          &iSBNEntityScrubber{},
		IP:            &iPEntityScrubber{},
		IBAN:          &iBANEntityScrubber{},
		PoBox:         &poBoxEntityScrubber{},
		ZipCode:       &zipCodeEntityScrubber{},
		CreditCard:    &creditCardEntityScrubber{},
		Phone:         &phoneEntityScrubber{},
		StreetAddress: &streetAddressEntityScrubber{},
		SSN:           &sSNEntityScrubber{},
		Link:          &linkEntityScrubber{},
		NotKnownPort:  &notKnownPortEntityScrubber{},
		SHA1Hex:       &sHA1HexEntityScrubber{},
		Time:          &timeEntityScrubber{},
		Date:          &dateEntityScrubber{},
		SHA256Hex:     &sHA256HexEntityScrubber{},
		GUID:          &gUIDEntityScrubber{},
		Email:         &emailEntityScrubber{},
		BtcAddress:    &btcAddressEntityScrubber{},
		GitRepo:       &gitRepoEntityScrubber{},
		StrictLink:    &strictLinkEntityScrubber{},
	}
)

// EntityConfig ...
type EntityConfig struct {
	ReplaceWith          *string
	MaskWithChar         *rune
	UnmaskedSuffixOffset int
	UnmaskedPrefixOffset int
}

func (e *EntityConfig) isValid() error {
	if e.UnmaskedPrefixOffset != 0 || e.UnmaskedSuffixOffset != 0 {
		if e.MaskWithChar == nil {
			return fmt.Errorf("prefix and suffix property can only be specified with MaskWithChar property")
		}

		return nil
	}

	if e.ReplaceWith != nil {
		return nil
	}

	if e.MaskWithChar != nil {
		return nil
	}

	return fmt.Errorf("either ReplaceWith or MaskWithChar must be specified")
}

type intermediateResponse struct {
	index    []int
	scrubber EntityScrubber
	entity   Entity
}

type intermediateScrubbingResponse struct {
	index int
	text  string
}

func (s *scrubber) sortIntervals(intervals []*intermediateResponse) {
	// sort intervals in the increasing order
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].index[0] != intervals[j].index[0] {
			return intervals[i].index[0] < intervals[j].index[0]
		}
		return intervals[i].index[1] >= intervals[j].index[1]
	})
}

func (s *scrubber) getEntityMatches(entities []Entity, text string) ([]*intermediateResponse, error) {
	intervals := make([]*intermediateResponse, 0)
	for _, entity := range entities {

		entityScrubber := _defaultEntityScrubbers[entity]
		if scrubber, ok := s.userProvidedScrubbers[entity]; ok {
			entityScrubber = scrubber
		}

		matches := entityScrubber.Match(text)
		if (len(matches)) > 0 {
			for _, match := range matches {
				if len(match) != 2 || match[0] >= match[1] {
					return nil, ErrInvalidMatchIndices
				}
				intervals = append(intervals, &intermediateResponse{
					index:    match,
					scrubber: entityScrubber,
					entity:   entity,
				})
			}
		}
	}
	return intervals, nil
}

func (s *scrubber) ScrubTexts(texts []string) ([]string, error) {

	wp := goworker.NewWorkerPool(&goworker.WorkerPoolInput{WorkerCount: 4})
	wp.Start()

	futures := make([]*goworker.Future, 0, len(texts))
	for i, val := range texts {
		index := i
		text := val

		futures = append(futures, wp.Add(&goworker.Task{
			F: func() (interface{}, error) {

				// sort find all the intervals ...
				intervals, err := s.getEntityMatches(s.blacklistedEntities, text)
				if err != nil {
					return nil, err
				}
				s.sortIntervals(intervals)

				nonOverlapping := make([]*intermediateResponse, 0, len(intervals))

				if len(intervals) > 0 {
					nonOverlapping = append(nonOverlapping, intervals[0])
				}

				// make intervals non overlapping
				for i := 1; i < len(intervals); i++ {
					if intervals[i].index[0] <= intervals[i-1].index[1] {
						if intervals[i-1].index[1] >= intervals[i].index[1] {
							continue
						}
						intervals[i].index[0] = intervals[i-1].index[1] + 1
					}
					nonOverlapping = append(nonOverlapping, intervals[i])
				}

				// remove intervals for ignored entities
				ignoredIntervals, err := s.getEntityMatches(s.ignoredEntities, text)
				if err != nil {
					return nil, err
				}
				s.sortIntervals(ignoredIntervals)

				scrubbable := make([]*intermediateResponse, 0)
				i, j := 0, 0
				for ; i < len(nonOverlapping) && j < len(ignoredIntervals); j++ {
					for ; i < len(nonOverlapping) && nonOverlapping[i].index[1] < ignoredIntervals[j].index[0]; i++ {
						scrubbable = append(scrubbable, nonOverlapping[i])
					}
					for ; i < len(nonOverlapping) && nonOverlapping[i].index[0] <= ignoredIntervals[j].index[1]; i++ {
					}
				}
				scrubbable = append(scrubbable, nonOverlapping[i:]...)

				intervals = scrubbable

				intervalsIterator := 0
				scrubbedText := make([]byte, 0, len(text))
				textBytes := []byte(text)
				txtIterator := 0
				for txtIterator < len(textBytes) {
					if intervalsIterator < len(intervals) && txtIterator == intervals[intervalsIterator].index[0] {
						config := _defaultEntityConfigs[intervals[intervalsIterator].entity]
						if val, ok := s.config[intervals[intervalsIterator].entity]; ok {
							config = val
						}
						replacementBytes := intervals[intervalsIterator].scrubber.Mask(textBytes[intervals[intervalsIterator].index[0]:intervals[intervalsIterator].index[1]], *config)
						scrubbedText = append(scrubbedText, replacementBytes...)
						txtIterator = intervals[intervalsIterator].index[1]
						intervalsIterator++
						continue
					}

					scrubbedText = append(scrubbedText, textBytes[txtIterator])
					txtIterator++
				}

				return &intermediateScrubbingResponse{
					index: index,
					text:  string(scrubbedText),
				}, nil

			},
		}))
	}

	wp.Done()
	wp.WaitForCompletion()

	results := make([]*intermediateScrubbingResponse, 0, len(texts))

	for _, future := range futures {
		fRes, fErr := future.Result(), future.Error()
		if fErr != nil {
			return nil, fErr
		}

		res := fRes.(*intermediateScrubbingResponse)
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	scrubbedTexts := make([]string, 0, len(texts))

	for _, val := range results {
		scrubbedTexts = append(scrubbedTexts, val.text)
	}

	return scrubbedTexts, nil
}

func (s *scrubber) ScrubStruct(obj interface{}) (interface{}, error) {
	return s.parse(obj)
}