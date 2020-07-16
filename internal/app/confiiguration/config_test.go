package confiiguration_test

import (
	"flag"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/PhilippZhulev/delta/internal/app/confiiguration"
	"github.com/fatih/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	var _ = Describe("Test Config", func() {
		var (
			configPath string
		)
		flag.StringVar(&configPath, "config-path", "configs/delta.toml", "path to config file")
		flag.Parse()
		config := confiiguration.NewConfig()
		_, err := toml.DecodeFile(configPath, config)
		n := structs.Names(config)
		v := structs.Values(config)

		It("valid cases, struct and not error", func() {
			Expect(err).To(HaveOccurred())
			Expect(len(n)).To(Equal(5))
			Expect(n[0]).To(Equal("BindAddr"))
			Expect(n[1]).To(Equal("LogLevel"))
			Expect(n[2]).To(Equal("DatabaseURL"))
			Expect(n[3]).To(Equal("SessionKey"))
			Expect(n[4]).To(Equal("Salt"))
		})

		It("valid cases no zero", func() {
			Expect(v[0]).ToNot(BeZero())
			Expect(v[1]).ToNot(BeZero())
		})

		It("valid cases no empty", func() {
			Expect(len(config.BindAddr)).ToNot(Equal(0))
			Expect(len(config.LogLevel)).ToNot(Equal(0))
		})
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}
