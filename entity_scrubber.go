package piiscrubber

// EntityScrubber ...
type EntityScrubber interface {
	Match(text string) [][]int
	Mask(detectedEntity []byte, config *EntityConfig) []byte
}

type mACAddressEntityScrubber struct {
}

func (s *mACAddressEntityScrubber) Match(text string) [][]int {
	return MACAddressRegex.FindAllStringIndex(text, -1)
}

func (s *mACAddressEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type mD5HexEntityScrubber struct {
}

func (s *mD5HexEntityScrubber) Match(text string) [][]int {
	return MD5HexRegex.FindAllStringIndex(text, -1)
}

func (s *mD5HexEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type iSBNEntityScrubber struct {
}

func (s *iSBNEntityScrubber) Match(text string) [][]int {
	indexes := ISBN10Regex.FindAllStringIndex(text, -1)
	indexes = append(indexes, ISBN13Regex.FindAllStringIndex(text, -1)...)
	return indexes
}

func (s *iSBNEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type iPEntityScrubber struct {
}

func (s *iPEntityScrubber) Match(text string) [][]int {
	return IPRegex.FindAllStringIndex(text, -1)
}

func (s *iPEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type iBANEntityScrubber struct {
}

func (s *iBANEntityScrubber) Match(text string) [][]int {
	return IBANRegex.FindAllStringIndex(text, -1)
}

func (s *iBANEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type poBoxEntityScrubber struct {
}

func (s *poBoxEntityScrubber) Match(text string) [][]int {
	return PoBoxRegex.FindAllStringIndex(text, -1)
}

func (s *poBoxEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type zipCodeEntityScrubber struct {
}

func (s *zipCodeEntityScrubber) Match(text string) [][]int {
	return ZipCodeRegex.FindAllStringIndex(text, -1)
}

func (s *zipCodeEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type creditCardEntityScrubber struct {
}

func (s *creditCardEntityScrubber) Match(text string) [][]int {
	return CreditCardRegex.FindAllStringIndex(text, -1)
}

func (s *creditCardEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type phoneEntityScrubber struct {
}

func (s *phoneEntityScrubber) Match(text string) [][]int {
	indexes := PhonesWithExtsRegex.FindAllStringIndex(text, -1)
	indexes = append(indexes, PhoneRegex.FindAllStringIndex(text, -1)...)

	return indexes
}

func (s *phoneEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type streetAddressEntityScrubber struct {
}

func (s *streetAddressEntityScrubber) Match(text string) [][]int {
	return StreetAddressRegex.FindAllStringIndex(text, -1)
}

func (s *streetAddressEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type sSNEntityScrubber struct {
}

func (s *sSNEntityScrubber) Match(text string) [][]int {
	return SSNRegex.FindAllStringIndex(text, -1)
}

func (s *sSNEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type linkEntityScrubber struct {
}

func (s *linkEntityScrubber) Match(text string) [][]int {
	return LinkRegex.FindAllStringIndex(text, -1)
}

func (s *linkEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type notKnownPortEntityScrubber struct {
}

func (s *notKnownPortEntityScrubber) Match(text string) [][]int {
	return NotKnownPortRegex.FindAllStringIndex(text, -1)
}

func (s *notKnownPortEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type sHA1HexEntityScrubber struct {
}

func (s *sHA1HexEntityScrubber) Match(text string) [][]int {
	return SHA1HexRegex.FindAllStringIndex(text, -1)
}

func (s *sHA1HexEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type timeEntityScrubber struct {
}

func (s *timeEntityScrubber) Match(text string) [][]int {
	return TimeRegex.FindAllStringIndex(text, -1)
}

func (s *timeEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type dateEntityScrubber struct {
}

func (s *dateEntityScrubber) Match(text string) [][]int {
	return DateRegex.FindAllStringIndex(text, -1)
}

func (s *dateEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type sHA256HexEntityScrubber struct {
}

func (s *sHA256HexEntityScrubber) Match(text string) [][]int {
	return SHA256HexRegex.FindAllStringIndex(text, -1)
}

func (s *sHA256HexEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type gUIDEntityScrubber struct {
}

func (s *gUIDEntityScrubber) Match(text string) [][]int {
	return GUIDRegex.FindAllStringIndex(text, -1)
}

func (s *gUIDEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type emailEntityScrubber struct {
}

func (s *emailEntityScrubber) Match(text string) [][]int {
	return EmailRegex.FindAllStringIndex(text, -1)
}

func (s *emailEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type btcAddressEntityScrubber struct {
}

func (s *btcAddressEntityScrubber) Match(text string) [][]int {
	return BtcAddressRegex.FindAllStringIndex(text, -1)
}

func (s *btcAddressEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type gitRepoEntityScrubber struct {
}

func (s *gitRepoEntityScrubber) Match(text string) [][]int {
	return GitRepoRegex.FindAllStringIndex(text, -1)
}

func (s *gitRepoEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

type strictLinkEntityScrubber struct {
}

func (s *strictLinkEntityScrubber) Match(text string) [][]int {
	return StrictLinkRegex.FindAllStringIndex(text, -1)
}

func (s *strictLinkEntityScrubber) Mask(detectedEntity []byte, config *EntityConfig) []byte {
	return NativeMasking(detectedEntity, config)
}

// NativeMasking ...
func NativeMasking(detectedEntity []byte, config *EntityConfig) []byte {
	if config.ReplaceWith != nil {
		return []byte(*config.ReplaceWith)
	}

	for index := config.UnmaskedPrefixOffset; index < len(detectedEntity)-config.UnmaskedSuffixOffset; index++ {
		detectedEntity[index] = byte(*config.MaskWithChar)
	}

	return detectedEntity
}
