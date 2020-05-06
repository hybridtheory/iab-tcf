package iab_tcf

import (
	"github.com/LiveRamp/iabconsent"
)

// ConsentV1 is an implementation of the Consent interface used to retrieve
// consent information given a TCF 1.0 format.
type ConsentV1 struct {
	ParsedConsent *iabconsent.ParsedConsent
}

// NewConsentV1 returns a consent interface from the reader received, with
// an error if something went wrong.
func NewConsentV1(reader *iabconsent.ConsentReader) (Consent, error) {
	parsedConsent, err := ParseV1(reader)
	if err != nil {
		return nil, err
	}
	return &ConsentV1{
		ParsedConsent: parsedConsent,
	}, nil
}

// Version returns the version of this consent string.
func (c *ConsentV1) Version() int {
	return int(iabconsent.V1)
}

// HasConsentedPurpose returns always true because consent TFC 1.0 doesn't
// come with this information.
func (c *ConsentV1) HasConsentedPurpose(purposeID int) bool {
	return c.ParsedConsent.PurposesAllowed[purposeID]
}

// GetConsentPurposeBitstring returns a string of 1 & 0 each of them representing the consent
// given for a specific purposeID.
func (c *ConsentV1) GetConsentPurposeBitstring() string {
	bitString := ""
	for i := 1; i <= 24; i++ {
		bitString += booleanFormatter[c.HasConsentedPurpose(i)]
	}
	return bitString
}

// HasConsentedLegitimateInterestForPurpose returns always true because consent TFC 1.0 doesn't
// come with this information.
func (c *ConsentV1) HasConsentedLegitimateInterestForPurpose(purposeID int) bool {
	return true
}

// HasUserConsented returns true if the user has given consent to the vendorID passed
// as parameter.
func (c *ConsentV1) HasUserConsented(vendorID int) bool {
	return c.ParsedConsent.VendorAllowed(vendorID)
}

// HasUserLegitimateInterest returns always true because consent TFC 1.0 doesn't
// come with this information.
func (c *ConsentV1) HasUserLegitimateInterest(vendorID int) bool {
	return true
}

// GetConsentBitstring returns a string of 1 & 0 each of them representing the consent
// given for a specific vendorID (the first number is for the vendorID 1, and so on).
func (c *ConsentV1) GetConsentBitstring() string {
	bitString := ""
	for i := 1; i <= c.ParsedConsent.MaxVendorID; i++ {
		bitString += booleanFormatter[c.HasUserConsented(i)]
	}
	return bitString
}

// GetInterestsBitstring returns an empty string always because consent TFC 1.0 doesn't
// implement user legitimate interests.
func (c *ConsentV1) GetInterestsBitstring() string {
	return ""
}

// GetPublisherRestrictions returns an empty list because consent TFC 1.0 doesn't
// implement user legitimate interests.
func (c *ConsentV1) GetPublisherRestrictions() ([]*iabconsent.PubRestrictionEntry) {
	return make([]*iabconsent.PubRestrictionEntry, 0, 0)
}

// ParseV1 uses a consent reader to extract information from a TCF 1.0 version
// consent string.
func ParseV1(r *iabconsent.ConsentReader) (*iabconsent.ParsedConsent, error) {
	var p = &iabconsent.ParsedConsent{}
	p.Version = int(iabconsent.V1)
	r.ReadString(17)
	p.ConsentLanguage, _ = r.ReadString(2)
	p.VendorListVersion, _ = r.ReadInt(12)
	p.PurposesAllowed, _ = r.ReadBitField(24)
	p.MaxVendorID, _ = r.ReadInt(16)
	p.IsRangeEncoding, _ = r.ReadBool()
	if p.IsRangeEncoding {
		p.DefaultConsent, _ = r.ReadBool()
		p.NumEntries, _ = r.ReadInt(12)
		p.RangeEntries, _ = r.ReadRangeEntries(uint(p.NumEntries))
	} else {
		p.ConsentedVendors, _ = r.ReadBitField(uint(p.MaxVendorID))
	}
	return p, r.Err
}
