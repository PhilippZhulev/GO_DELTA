package validate_test

import (
	"testing"

	"github.com/PhilippZhulev/delta/internal/app/validate"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Тестовая структура
type test struct {
	Field string
}

func TestValidation(t *testing.T) {

	var _ = Describe("Test validation password", func() {
		It("no valid cases", func() {
			Expect(validate.Pass("abvsgdez")).To(HaveOccurred())
			Expect(validate.Pass("ZBSDWAC")).To(HaveOccurred())
			Expect(validate.Pass("12345678")).To(HaveOccurred())
			Expect(validate.Pass("abc")).To(HaveOccurred())
			Expect(validate.Pass("a1bcZ2")).To(HaveOccurred())
			Expect(validate.Pass("P@ssword")).To(HaveOccurred())
		})

		It("valid cases", func() {
			Expect(validate.Pass("AbcZ76vn")).To(BeNil())
			Expect(validate.Pass("1Hb673Va")).To(BeNil())
			Expect(validate.Pass("d@Bv2Nv23Ba")).To(BeNil())
		})
	})

	var _ = Describe("Test confirm", func() {
		It("valid confirm cases", func() {
			Expect(validate.Confirm("Password1", "Password1")).To(BeNil())
		})
		It("no valid confirm cases", func() {
			Expect(validate.Confirm("Password1", "Password2")).To(HaveOccurred())
		})
	})

	var _ = Describe("Test requiredIf", func() {
		t := &test{}

		It("valid requiredIf cases", func() {
			Expect(validate.RequiredIf(true)(t.Field)).To(HaveOccurred())
		})
		It("no valid requiredIf cases", func() {
			t.Field = "test"
			Expect(validate.RequiredIf(true)(t.Field)).NotTo(HaveOccurred())
		})
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}