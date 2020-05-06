package iab_tcf

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Consent TCF 1.0", func() {

	const (
		testGdprConsent = "BOlLbqtOlLbqtAVABADECg-AAAApp7v______9______9uz_Ov_v_f__33e8__9v_l_7_-___u_-3zd4u_1vf99yfm1-7etr3tp_87ues2_Xur__79__3z3_9phP78k89r7337Ew-v02"
		testBitstringConsentPurpose = "111110000000000000000000"
		testBitstringConsent = "111101110111111111111111111111111111111111111111111110111111111111111111111111111111111111111110110111011001111111100111010111111111110111111111101111111111111111111011111011101111011110011111111111111110110111111111110010111111111101111111111111011111111111111111110111011111111111011011111001101110111100010111011111111010110111101111111110111110111001001111110011011010111111011101101111010110110101111011110110110100111111111110011101110111001111010110011011011111101011110111010101111111111111111101111110111111111111111011111001111011111111111110110100110000100111111101111110010010011110011110110101111101111011111011111101100010011000011111010111111010011011"
	)

	var (
		consent Consent
		err error
	)

	BeforeEach(func() {
		consent, err = NewConsent(testGdprConsent)
		Expect(err).NotTo(HaveOccurred())		
	})

	It("detects the version as V1", func() {
		Expect(consent.Version()).To(Equal(1))
	})

	DescribeTable("purpose consent",
		func(purposeID int, expected bool) {
			Expect(consent.HasConsentedPurpose(purposeID)).To(Equal(expected))
		},
		Entry("the purpose id -1", -1, false),
		Entry("the purpose id 0", 0, false),
		Entry("the purpose id 1", 1, true),
		Entry("the purpose id 5", 5, true),
		Entry("the purpose id 20", 20, false),
		Entry("the purpose id 10000", 10000, false),
	)

	It("returns the purpose consent bitstring", func() {
		Expect(consent.GetConsentPurposeBitstring()).To(Equal(testBitstringConsentPurpose))
	})

	It("returns always true for 'HasConsentedLegitimateInterestForPurpose' because TCF 1.0 does not contain that info", func() {
		for purposeID := 1; purposeID <= 24; purposeID++ {
			Expect(consent.HasConsentedLegitimateInterestForPurpose(purposeID)).To(BeTrue())
		}
	})

	DescribeTable("vendors allowed",
		func(vendorID int, expected bool) {
			Expect(consent.HasUserConsented(vendorID)).To(Equal(expected))
		},
		Entry("the vendor id -1", -1, false),
		Entry("the vendor id 0", 0, false),
		Entry("the vendor id 1", 1, true),
		Entry("the vendor id 2", 2, true),
		Entry("the vendor id 50", 50, true),
		Entry("the vendor id 99", 99, false),
		Entry("the vendor id 150", 150, true),
		Entry("the vendor id 204", 204, false),
		Entry("the vendor id 250", 250, true),
		Entry("the vendor id 300", 300, false),
		Entry("the vendor id 665", 665, true),
		Entry("the vendor id 665", 666, true),
		Entry("the vendor id 665", 667, false),
		Entry("the vendor id 10000", 10000, false),
	)

	It("returns always true for HasUserLegitimateInterest because TCF 1.0 does not contain that info", func() {
		for vendorID := 1; vendorID <= 1000; vendorID++ {
			Expect(consent.HasUserLegitimateInterest(vendorID)).To(BeTrue())
		}
	})

	It("returns the consent bitstring", func() {
		Expect(consent.GetConsentBitstring()).To(Equal(testBitstringConsent))
	})

	It("always returns the user interest bitstring as empty because TCF 1.0 does not contain that info", func() {
		Expect(consent.GetInterestsBitstring()).To(Equal(""))
	})

	It("always returns the publisher restrictions empty because TCF 1.0 does not contain that info", func() {
		Expect(consent.GetPublisherRestrictions()).To(HaveLen(0))
	})
})
