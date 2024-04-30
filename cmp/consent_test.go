package cmp_test

import (
	"github.com/hybridtheory/iab-tcf/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consent", func() {
	var (
		testValidCMPs = []int{1, 2, 3}
		consent       *cmp.Consent
	)

	BeforeEach(func() {
		consent = &cmp.Consent{}
	})

	It("returns valid cmps", func() {
		cmp.ValidCMPs = testValidCMPs
		Expect(consent.ValidCMPs()).To(Equal(testValidCMPs))
	})

	Context("is list loaded", func() {
		It("returns false if it is not loaded", func() {
			cmp.ValidCMPs = nil
			Expect(consent.IsCMPListLoaded()).To(BeFalse())
		})

		It("returns true if it is loaded but empty", func() {
			cmp.ValidCMPs = []int{}
			Expect(consent.IsCMPListLoaded()).To(BeTrue())
		})

		It("returns true if it is properly loaded", func() {
			cmp.ValidCMPs = testValidCMPs
			Expect(consent.IsCMPListLoaded()).To(BeTrue())
		})
	})
})
