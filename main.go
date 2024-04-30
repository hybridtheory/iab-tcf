package iab_tcf

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/LiveRamp/iabconsent"
)

var booleanFormatter = map[bool]string{
	true:  "1",
	false: "0",
}

// Consent is an interface to retrieve the most important information
// from consent strings, no matter if they are v1 or v2.
type Consent interface {
	// Version returns the version of this consent string.
	Version() int
	// CMPID returns the CMP ID of this consent string.
	CMPID() int
	// HasConsentedPurpose returns the consent value for a Purpose established on the legal basis of consent.
	// The Purposes are numerically identified and published in the Global Vendor List.
	HasConsentedPurpose(purposeID int) bool
	// GetConsentPurposeBitstring returns a string of 1 & 0 each of them representing the consent
	// given for a specific purposeID.
	// The first number is for the purposeID 1, the second number for the purposeID 2,
	// and so on.
	GetConsentPurposeBitstring() string
	// HasConsentedLegitimateInterestForPurpose returns the Purpose’s transparency requirements
	// are met for each Purpose on the legal basis of legitimate interest and the user has not
	// exercised their “Right to Object” to that Purpose.
	HasConsentedLegitimateInterestForPurpose(purposeID int) bool
	// HasUserConsented returns true if the user has given consent to the vendorID passed
	// as parameter.
	HasUserConsented(vendorID int) bool
	// HasUserLegitimateInterest returns true if the CMP has established transparency for a vendor's
	// legitimate interest disclosures. If a user exercises their “Right To Object” to a vendor’s
	// processing based on a legitimate interest, then it returns false.
	HasUserLegitimateInterest(vendorID int) bool
	// GetConsentBitstring returns a string of 1 & 0 each of them representing the consent
	// given for a specific vendorID.
	// The first number is for the vendorID 1, the second number for the vendorID 2,
	// and so on.
	GetConsentBitstring() string
	// GetInterestsBitstring returns a string of 1 & 0 each of them representing the
	// return of `HasUserLegitimateInterest` method for the vendorID in that position of the string.
	// The first number is for the vendorID 1, the second number for the vendorID 2,
	// and so on.
	GetInterestsBitstring() string
	// GetPublisherRestrictions returns a list of restrictions per publisher, if it relates.
	GetPublisherRestrictions() []*iabconsent.PubRestrictionEntry
	// IsCMPListLoaded returns if the list of valid CMPs was properly loaded or not.
	IsCMPListLoaded() bool
	// IsCMPValid validates the consent string CMP ID agains the list of valid ones downloaded from IAB.
	IsCMPValid() bool
}

// DecodeConsent receives a GDPR IAB consent string and decodes the
// CORE segment only, returning it. It also returns an error if something
// happened and we couldn't decode it.
func DecodeConsent(consent string) ([]byte, error) {
	segments := strings.Split(consent, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(segments[0])
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// GetVersion extracts the version from the consent string, moving
// the data pointer so we don't have to reparse it.
func GetVersion(r *iabconsent.ConsentReader) iabconsent.TCFVersion {
	if version, err := r.ReadInt(6); err == nil {
		return iabconsent.TCFVersion(version)
	}
	return iabconsent.InvalidTCFVersion
}

// NewConsent returns a Consent instance with all the necessary information
// available. It returns an error if something went wrong.
func NewConsent(consent string) (Consent, error) {
	decoded, err := DecodeConsent(consent)
	if err == nil {
		reader := iabconsent.NewConsentReader(decoded)
		switch GetVersion(reader) {
		case iabconsent.V1:
			return NewConsentV1(reader)
		case iabconsent.V2:
			return NewConsentV2(reader)
		}
		return nil, errors.New("Invalid consent version found")
	}
	return nil, err
}
