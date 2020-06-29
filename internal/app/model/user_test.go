package model_test

import (
	"testing"

	"github.com/PhilippZhulev/delta/internal/app/model"
	"github.com/fatih/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


func TestUserModel(t *testing.T) {
	us := &model.User{}

	var _ = Describe("Structure test", func() {
		It("valid cases, struct len", func() {
			n := structs.Names(us)
			Expect(len(n)).To(Equal(11))
		})

		It("valid cases, struct names", func() {
			n := structs.Names(us)
			Expect(n[0]).To(Equal("ID"))
			Expect(n[1]).To(Equal("Login"))
			Expect(n[2]).To(Equal("Password"))
			Expect(n[3]).To(Equal("EncryptedPassword"))
			Expect(n[4]).To(Equal("СonfirmEncryptedPassword"))
			Expect(n[5]).To(Equal("JobCode"))
			Expect(n[6]).To(Equal("Email"))
			Expect(n[7]).To(Equal("Phone"))
			Expect(n[8]).To(Equal("Name"))
			Expect(n[9]).To(Equal("UUID"))
			Expect(n[10]).To(Equal("Role"))
		})
	})

	var _ = Describe("Test Validate", func() {
		It("valid cases", func() {
			us.Login = "Tester-US" 
			us.EncryptedPassword = "1Sd8C23b"
			Expect(us.Validate()).To(BeNil())
		})

		It("not valid cases, invaalid name", func() {
			us.Login = "Tes" 
			us.EncryptedPassword = "1Sd8C23b"
			Expect(us.Validate()).NotTo(BeNil())
		})

		It("not valid cases, invaalid pass", func() {
			us.Login = "Tester-US" 
			us.EncryptedPassword = "1S23456"
			Expect(us.Validate()).NotTo(BeNil())
		})

		It("not valid cases, invaalid pass and invalid name", func() {
			us.Login = "Tes"
			us.EncryptedPassword = "1S23456"
			Expect(us.Validate()).NotTo(BeNil())
		})
	})

	var _ = Describe("Test ValidatePassword", func() {
		It("valid cases", func() {
			Expect(us.ValidatePassword("1Sd8C23b", "1Sd8C23b")).To(BeNil())
		})

		It("no valid cases, different, first", func() {
			Expect(us.ValidatePassword("123", "1Sd8C23b")).NotTo(BeNil())
		})

		It("no valid cases, different, first empty", func() {
			Expect(us.ValidatePassword("", "1Sd8C23b")).NotTo(BeNil())
		})

		It("no valid cases, different, first and last empty", func() {
			Expect(us.ValidatePassword("", "")).NotTo(BeNil())
		})
	})


	var _ = Describe("Test Sanitize", func() {
		It("valid cases", func() {
			us.Password = "1Sd8C23b"
			us.EncryptedPassword = "KNVDVSVLJD"
			us.Sanitize()
			Expect(us.Password).To(Equal(""))
			Expect(us.EncryptedPassword).To(Equal(""))
		})

		It("not valid cases", func() {
			us.Password = "1Sd8C23b"
			us.EncryptedPassword = "KNVDVSVLJD"
			us.Sanitize()
			Expect(us.Password).NotTo(Equal("1Sd8C23b"))
			Expect(us.EncryptedPassword).NotTo(Equal("KNVDVSVLJD"))
		})
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}