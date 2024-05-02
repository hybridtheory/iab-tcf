package iab_tcf_test

import (
	iab_tcf "github.com/hybridtheory/iab-tcf"
	"github.com/hybridtheory/iab-tcf/cmp"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	. "github.com/onsi/gomega"
)

var _ = Describe("Consent TCF 2.0", func() {

	const (
		testGdprConsent              = "COxR03kOxR1CqBcABCENAgCMAP_AAH_AAAqIF3EXySoGY2thI2YVFxBEIYwfJxyigMgChgQIsSwNQIeFLBoGLiAAHBGYJAQAGBAEEACBAQIkHGBMCQAAgAgBiRCMQEGMCzNIBIBAggEbY0FACCVmHkHSmZCY7064O__QLuIJEFQMAkSBAIACLECIQwAQDiAAAYAlAAABAhIaAAgIWBQEeAAAACAwAAgAAABBAAACAAQAAICIAAABAAAgAiAQAAAAGgIQAACBABACRIAAAEANCAAgiCEAQg4EAo4AAA"
		testBitstringConsent         = "010001011111001001001010100000011001100011011010110110000100100011011001100001010100010111000100000100010000100001100011000001111100100111000111001010001010000000110010000000001010000110000001000000100010110001001011000000110101000000100001111000010100101100000110100000011000101110001000000000000000011100000100011001100000100100000001000000000000011000000100000000010000010000000000001000000100000001000000100010010000011100011000000100110000001001000000000000000010000000000010000000000110001001000100001000110001000000010000011000110000001011001100110100100000000100100000000100000010000010000000010001101101100011010000010100000000001000001001010110011000011110010000011101001010011001100100001001100011101111010011101011100000111011111111111101"
		testBitstringConsentPurpose  = "111111111100000000000000"
		testBitstringInterestConsent = "010000010010001000001010100000011000000001001000100100000010000000010000000000000100010110001000000100010000100001100000000000100000000111000100000000000000000000110000000001001010000000000000000000000010000001000010010000110100000000000001000000010000101100000010100000001000111100000000000000000000000000000100000001100000000000000001000000000000000000000000000010000010000000000000000000000100000000000000100000000000000000010000000100010000000000000000000000000010000000000000000001000000000001000100000000100000000000000000000000000000001101000000010000100000000000000000000100000010000000000010000000000100100010010000000000000000000000001000000000011010000100000000000001000001000100000100001000000000100001000001110000001000000001010001110000"
	)

	var (
		consent iab_tcf.Consent
		err     error
	)

	BeforeEach(func() {
		consent, err = iab_tcf.NewConsent(testGdprConsent)
		Expect(err).NotTo(HaveOccurred())
	})

	It("detects the version as V2", func() {
		Expect(consent.Version()).To(Equal(2))
	})

	It("detects the CMP ID as 21", func() {
		Expect(consent.CMPID()).To(Equal(92))
	})

	DescribeTable("purposes consented",
		func(purposeID int, expected bool) {
			Expect(consent.HasConsentedPurpose(purposeID)).To(Equal(expected))
		},
		Entry("the purpose id -1", -1, false),
		Entry("the purpose id 0", 0, false),
		Entry("the purpose id 2", 2, true),
		Entry("the purpose id 4", 4, true),
		Entry("the purpose id 6", 6, true),
		Entry("the purpose id 8", 8, true),
		Entry("the purpose id 10", 10, true),
		Entry("the purpose id 11", 11, false),
		Entry("the purpose id 10000", 10000, false),
	)

	It("returns the purpose consent bitstring", func() {
		Expect(consent.GetConsentPurposeBitstring()).To(Equal(testBitstringConsentPurpose))
	})

	DescribeTable("legitimate interest for purposes",
		func(purposeID int, expected bool) {
			Expect(consent.HasConsentedLegitimateInterestForPurpose(purposeID)).To(Equal(expected))
		},
		Entry("the purpose id -1", -1, false),
		Entry("the purpose id 1", 1, false),
		Entry("the purpose id 2", 2, true),
		Entry("the purpose id 8", 8, true),
		Entry("the purpose id 9", 9, true),
		Entry("the purpose id 10", 10, true),
		Entry("the purpose id 11", 11, false),
		Entry("the purpose id 10000", 10000, false),
	)

	DescribeTable("vendors allowed",
		func(vendorID int, expected bool) {
			Expect(consent.HasUserConsented(vendorID)).To(Equal(expected))
		},
		Entry("the vendor id -1", -1, false),
		Entry("the vendor id 0", 0, false),
		Entry("the vendor id 1", 1, false),
		Entry("the vendor id 2", 2, true),
		Entry("the vendor id 50", 50, true),
		Entry("the vendor id 99", 99, false),
		Entry("the vendor id 150", 150, false),
		Entry("the vendor id 204", 204, false),
		Entry("the vendor id 250", 250, true),
		Entry("the vendor id 300", 300, false),
		Entry("the vendor id 665", 665, true),
		Entry("the vendor id 666", 666, false),
		Entry("the vendor id 667", 667, false),
		Entry("the vendor id 750", 750, true),
		Entry("the vendor id 751", 751, false),
		Entry("the vendor id 10000", 10000, false),
	)

	DescribeTable("vendors legitimate interests",
		func(vendorID int, expected bool) {
			Expect(consent.HasUserLegitimateInterest(vendorID)).To(Equal(expected))
		},
		Entry("the vendor id -1", -1, false),
		Entry("the vendor id 0", 0, false),
		Entry("the vendor id 1", 1, false),
		Entry("the vendor id 2", 2, true),
		Entry("the vendor id 50", 50, false),
		Entry("the vendor id 99", 99, false),
		Entry("the vendor id 150", 150, false),
		Entry("the vendor id 204", 204, false),
		Entry("the vendor id 250", 250, false),
		Entry("the vendor id 300", 300, false),
		Entry("the vendor id 665", 665, false),
		Entry("the vendor id 665", 666, false),
		Entry("the vendor id 665", 667, false),
		Entry("the vendor id 746", 746, true),
		Entry("the vendor id 747", 747, false),
		Entry("the vendor id 10000", 10000, false),
	)

	It("returns the consent bitstring", func() {
		Expect(consent.GetConsentBitstring()).To(Equal(testBitstringConsent))
	})

	It("returns the user interest bitstring", func() {
		Expect(consent.GetInterestsBitstring()).To(Equal(testBitstringInterestConsent))
	})

	It("returns the publisher restrictions", func() {
		Expect(consent.GetPublisherRestrictions()).To(HaveLen(0))
	})
})

var _ = Describe("Consent TCF 2.0 generated with https://iabtcf.com/#/encode", func() {

	const (
		testGdprConsent              = "COytyllOytyllCrAAAENAiCMAFVAACqAAAAAF3QAgAFABkAAoioAAA.IF5EX2S5OI2tho2YdF7BEYYwfJxyigMgShgQIsS8NwIeFbBoGPmAAHBG4JAQAGBAkkACBAQIsHGBcCQABgIgRiRCMQEGMjzNKBJBAggkbI0FACCVmnkHS3ZCY70-6u__bA"
		testBitstringConsent         = "000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
		testBitstringConsentPurpose  = "010101010100000000000000"
		testBitstringInterestConsent = "0100010101"
	)

	var (
		consent iab_tcf.Consent
		err     error
	)

	BeforeEach(func() {
		consent, err = iab_tcf.NewConsent(testGdprConsent)
		Expect(err).NotTo(HaveOccurred())
	})

	It("detects the version as V2", func() {
		Expect(consent.Version()).To(Equal(2))
	})

	It("detects the cmp id as 171", func() {
		Expect(consent.CMPID()).To(Equal(171))
	})

	DescribeTable("cmp validity",
		func(validCMPs []int, expected gomega.OmegaMatcher) {
			cmp.ValidCMPs = validCMPs
			Expect(consent.IsCMPValid()).To(expected)
		},
		Entry("valid", []int{171}, BeTrue()),
		Entry("invalid", []int{172}, BeFalse()),
	)

	DescribeTable("purposes consented",
		func(purposeID int, expected bool) {
			Expect(consent.HasConsentedPurpose(purposeID)).To(Equal(expected))
		},
		Entry("the purpose id -1", -1, false),
		Entry("the purpose id 0", 0, false),
		Entry("the purpose id 2", 2, true),
		Entry("the purpose id 4", 4, true),
		Entry("the purpose id 6", 6, true),
		Entry("the purpose id 8", 8, true),
		Entry("the purpose id 10", 10, true),
		Entry("the purpose id 11", 11, false),
		Entry("the purpose id 10000", 10000, false),
	)

	It("returns the purpose consent bitstring", func() {
		Expect(consent.GetConsentPurposeBitstring()).To(Equal(testBitstringConsentPurpose))
	})

	DescribeTable("legitimate interest for purposes",
		func(purposeID int, expected bool) {
			Expect(consent.HasConsentedLegitimateInterestForPurpose(purposeID)).To(Equal(expected))
		},
		Entry("the purpose id -1", -1, false),
		Entry("the purpose id 1", 1, false),
		Entry("the purpose id 3", 3, true),
		Entry("the purpose id 5", 5, true),
		Entry("the purpose id 7", 7, true),
		Entry("the purpose id 9", 9, true),
		Entry("the purpose id 10", 10, false),
		Entry("the purpose id 11", 11, false),
		Entry("the purpose id 10000", 10000, false),
	)

	DescribeTable("vendors allowed",
		func(vendorID int, expected bool) {
			Expect(consent.HasUserConsented(vendorID)).To(Equal(expected))
		},
		Entry("the vendor id 0", 0, false),
		Entry("the vendor id 1", 1, false),
		Entry("the vendor id 10", 10, true),
		Entry("the vendor id 100", 100, true),
		Entry("the vendor id 150", 150, false),
		Entry("the vendor id 10000", 10000, false),
	)

	DescribeTable("vendors legitimate interests",
		func(vendorID int, expected bool) {
			Expect(consent.HasUserLegitimateInterest(vendorID)).To(Equal(expected))
		},
		Entry("the vendor id 0", 0, false),
		Entry("the vendor id 2", 2, true),
		Entry("the vendor id 4", 4, false),
		Entry("the vendor id 6", 6, true),
		Entry("the vendor id 8", 8, true),
		Entry("the vendor id 10", 10, true),
	)

	It("returns the consent bitstring", func() {
		Expect(consent.GetConsentBitstring()).To(Equal(testBitstringConsent))
	})

	It("returns the user interest bitstring", func() {
		Expect(consent.GetInterestsBitstring()).To(Equal(testBitstringInterestConsent))
	})

	It("returns the publisher restrictions", func() {
		Expect(consent.GetPublisherRestrictions()).To(HaveLen(0))
	})
})
