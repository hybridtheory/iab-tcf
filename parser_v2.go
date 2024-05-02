package iab_tcf

import (
	"slices"

	"github.com/LiveRamp/iabconsent"
	"github.com/hybridtheory/iab-tcf/cmp"
)

// ConsentV2 is an implementation of the Consent interface used to retrieve
// consent information given a TCF 1.0 format.
type ConsentV2 struct {
	cmp.Consent
	ParsedConsent *iabconsent.V2ParsedConsent
}

// NewConsentV2 returns a consent interface from the reader received, with
// an error if something went wrong.
func NewConsentV2(reader *iabconsent.ConsentReader) (Consent, error) {
	parsedConsent, err := ParseV2(reader)
	if err != nil {
		return nil, err
	}
	return &ConsentV2{
		ParsedConsent: parsedConsent,
	}, nil
}

// Version returns the version of this consent string.
func (c *ConsentV2) Version() int {
	return int(iabconsent.V2)
}

// CMPID returns the CMP ID of this consent string.
func (c *ConsentV2) CMPID() int {
	return c.ParsedConsent.CMPID
}

// IsCMPValid validates the consent string CMP ID agains the list of valid ones downloaded from IAB.
func (c *ConsentV2) IsCMPValid() bool {
	return slices.Contains(c.ValidCMPs(), c.CMPID())
}

// HasConsentedPurpose returns the consent value for a Purpose established on the legal basis of consent.
// The Purposes are numerically identified and published in the Global Vendor List.
func (c *ConsentV2) HasConsentedPurpose(purposeID int) bool {
	return c.ParsedConsent.PurposesConsent[purposeID]
}

// GetConsentPurposeBitstring returns a string of 1 & 0 each of them representing the consent
// given for a specific purposeID.
func (c *ConsentV2) GetConsentPurposeBitstring() string {
	bitString := ""
	for i := 1; i <= 24; i++ {
		bitString += booleanFormatter[c.HasConsentedPurpose(i)]
	}
	return bitString
}

// HasConsentedLegitimateInterestForPurpose returns the Purpose’s transparency requirements
// are met for each Purpose on the legal basis of legitimate interest and the user has not
// exercised their “Right to Object” to that Purpose.
func (c *ConsentV2) HasConsentedLegitimateInterestForPurpose(purposeID int) bool {
	return c.ParsedConsent.PurposesLITransparency[purposeID]
}

// HasUserConsented returns true if the user has given consent to the vendorID passed
// as parameter.
func (c *ConsentV2) HasUserConsented(vendorID int) bool {
	if c.ParsedConsent.IsConsentRangeEncoding {
		for _, re := range c.ParsedConsent.ConsentedVendorsRange {
			if re.StartVendorID <= vendorID && vendorID <= re.EndVendorID {
				return true
			}
		}
		return false
	}
	return c.ParsedConsent.ConsentedVendors[vendorID]
}

// HasUserLegitimateInterest returns true if the CMP has established transparency for a vendor's
// legitimate interest disclosures. If a user exercises their “Right To Object” to a vendor’s
// processing based on a legitimate interest, then it returns false.
func (c *ConsentV2) HasUserLegitimateInterest(vendorID int) bool {
	if c.ParsedConsent.IsInterestsRangeEncoding {
		for _, re := range c.ParsedConsent.InterestsVendorsRange {
			if re.StartVendorID <= vendorID && vendorID <= re.EndVendorID {
				return true
			}
		}
		return false
	}
	return c.ParsedConsent.InterestsVendors[vendorID]
}

// GetConsentBitstring returns a string of 1 & 0 each of them representing the consent
// given for a specific vendorID (the first number is for the vendorID 1, and so on).
func (c *ConsentV2) GetConsentBitstring() string {
	bitString := ""
	for i := 1; i <= c.ParsedConsent.MaxConsentVendorID; i++ {
		bitString += booleanFormatter[c.HasUserConsented(i)]
	}
	return bitString
}

// GetInterestsBitstring returns a string of 1 & 0 each of them representing the user's interest to
// `Right To Object` for a specific vendorID (the first number is for the vendorID 1, and so on).
func (c *ConsentV2) GetInterestsBitstring() string {
	bitString := ""
	for i := 1; i <= c.ParsedConsent.MaxInterestsVendorID; i++ {
		bitString += booleanFormatter[c.HasUserLegitimateInterest(i)]
	}
	return bitString
}

// GetPublisherRestrictions returns a list of restrictions per publisher, if it relates.
func (c *ConsentV2) GetPublisherRestrictions() []*iabconsent.PubRestrictionEntry {
	return c.ParsedConsent.PubRestrictionEntries
}

// ParseV2 uses a consent reader to extract information from a TCF 2.0 version
// consent string.
func ParseV2(r *iabconsent.ConsentReader) (*iabconsent.V2ParsedConsent, error) {
	var p = &iabconsent.V2ParsedConsent{}
	p.Version = int(iabconsent.V2)
	r.ReadString(12)
	p.CMPID, _ = r.ReadInt(12)
	r.ReadString(3)
	p.ConsentLanguage, _ = r.ReadString(2)
	p.VendorListVersion, _ = r.ReadInt(12)
	p.TCFPolicyVersion, _ = r.ReadInt(6)
	p.IsServiceSpecific, _ = r.ReadBool()
	p.UseNonStandardStacks, _ = r.ReadBool()
	p.SpecialFeaturesOptIn, _ = r.ReadBitField(12)
	p.PurposesConsent, _ = r.ReadBitField(24)
	p.PurposesLITransparency, _ = r.ReadBitField(24)
	p.PurposeOneTreatment, _ = r.ReadBool()
	p.PublisherCC, _ = r.ReadString(2)
	p.MaxConsentVendorID, _ = r.ReadInt(16)
	p.IsConsentRangeEncoding, _ = r.ReadBool()
	if p.IsConsentRangeEncoding {
		p.NumConsentEntries, _ = r.ReadInt(12)
		p.ConsentedVendorsRange, _ = r.ReadRangeEntries(uint(p.NumConsentEntries))
	} else {
		p.ConsentedVendors, _ = r.ReadBitField(uint(p.MaxConsentVendorID))
	}
	p.MaxInterestsVendorID, _ = r.ReadInt(16)
	p.IsInterestsRangeEncoding, _ = r.ReadBool()
	if p.IsInterestsRangeEncoding {
		p.NumInterestsEntries, _ = r.ReadInt(12)
		p.InterestsVendorsRange, _ = r.ReadRangeEntries(uint(p.NumInterestsEntries))
	} else {
		p.InterestsVendors, _ = r.ReadBitField(uint(p.MaxInterestsVendorID))
	}
	p.NumPubRestrictions, _ = r.ReadInt(12)
	p.PubRestrictionEntries, _ = r.ReadPubRestrictionEntries(uint(p.NumPubRestrictions))
	return p, r.Err
}
