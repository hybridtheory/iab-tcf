package cmp_test

import (
	"os"

	"github.com/hybridtheory/iab-tcf/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Loader", func() {

	var (
		loader *cmp.Loader
		cmps   []cmp.CMP
		err    error
	)

	BeforeEach(func() {
		loader = cmp.NewLoader()
	})

	Describe("configuration", func() {
		Context("with url", func() {
			const (
				testURL = "https://unknown-url"
			)

			BeforeEach(func() {
				loader = cmp.NewLoader(cmp.WithURL(testURL))
				err = loader.LoadIDs()
			})

			It("triggers an error because the url is not found", func() {
				Expect(err).To(HaveOccurred())
			})

			It("is used to retrieve the JSON", func() {
				Expect(err).Should(MatchError(MatchRegexp("lookup unknown-url")))
			})
		})

		Context("with json", func() {
			const (
				testJSONFile = "cmp_test.json"
			)

			BeforeEach(func() {
				contents, _ := os.ReadFile(testJSONFile)
				loader = cmp.NewLoader(cmp.WithJSON(string(contents)))
				err = loader.LoadIDs()
			})

			It("does not trigger an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("is used to parse the JSON", func() {
				Expect(cmp.ValidCMPs).Should(Equal([]int{1, 2}))
			})
		})
	})

	Describe("load", func() {
		BeforeEach(func() {
			cmps, err = loader.Load()
		})

		DescribeTable("available vendors",
			func(cmpID int) {
				Expect(cmps).To(ContainElement(MatchFields(IgnoreExtras, Fields{"ID": Equal(cmpID)})))
			},
			Entry("Microsoft Corporation", 198),
			Entry("Google LLC", 300),
			Entry("eBay Inc", 125),
		)

		DescribeTable("unavailable vendors",
			func(cmpID int) {
				Expect(cmps).ToNot(ContainElement(MatchFields(IgnoreExtras, Fields{"ID": Equal(cmpID)})))
			},
			Entry("unknown #1", 4),
			Entry("unknown #2", 8),
		)

		Context("errors", func() {
			It("unavailable endpoint", func() {
				cmps, err = cmp.NewLoader(cmp.WithURL("https://unknown")).Load()
				Expect(cmps).To(HaveLen(0))
				Expect(err).To(HaveOccurred())
			})

			It("not a json", func() {
				cmps, err = cmp.NewLoader(cmp.WithURL("http://github.com/")).Load()
				Expect(cmps).To(HaveLen(0))
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("load ids", func() {
		BeforeEach(func() {
			err = loader.LoadIDs()
		})

		DescribeTable("available vendors",
			func(expected int) {
				Expect(cmp.ValidCMPs).To(ContainElement(expected))
			},
			Entry("Microsoft Corporation", 198),
			Entry("Google LLC", 300),
			Entry("eBay Inc", 125),
		)

		DescribeTable("unavailable vendors",
			func(unexpected int) {
				Expect(cmp.ValidCMPs).ToNot(ContainElement(unexpected))
			},
			Entry("unknown #1", 4),
			Entry("unknown #2", 8),
		)

		Context("errors", func() {
			BeforeEach(func() {
				cmp.ValidCMPs = nil
			})

			It("unavailable endpoint", func() {
				err = cmp.NewLoader(cmp.WithURL("https://unknown")).LoadIDs()
				Expect(cmps).To(HaveLen(0))
				Expect(err).To(HaveOccurred())
			})

			It("not a json", func() {
				err = cmp.NewLoader(cmp.WithURL("http://github.com/")).LoadIDs()
				Expect(cmp.ValidCMPs).To(HaveLen(0))
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
