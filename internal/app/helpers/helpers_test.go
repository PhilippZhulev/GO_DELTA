package helpers_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/PhilippZhulev/delta/internal/app/helpers"
	"github.com/PhilippZhulev/delta/internal/app/store/sqlstore"
	_ "github.com/lib/pq" // ...
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHelpers(t *testing.T) {

	var it struct{
		hesh helpers.Hesh
		store sqlstore.Store
		respond helpers.Respond
	}

	db, err := sql.Open("postgres", "host=localhost dbname=delta_test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	store := sqlstore.New(db)
	
	var _ = Describe("Test hesh password", func() {
		It("no valid cases", func() {
			Expect(it.hesh.HashPassword("123456")).NotTo(Equal("123456"))
			Expect(it.hesh.HashPassword("123456")).NotTo(Equal("test"))
		})

		It("valid cases", func() {
			Expect(it.hesh.HashPassword("123456")).To(Equal("29l4zNu+i23nf2s3td+bW2Kn6JKlAcO1PqoWsIOL1e0="))
			Expect(it.hesh.HashPassword("test")).To(Equal("Q7DO+ZJl+eNMEOqdNQGSbSezn1fG1nRWHYuiNueoGfs="))
			Expect(it.hesh.HashPassword("")).To(Equal("thNnmggU2ex3L5XXeMNfxf8Wl8STcVZTxscSFEKSxa0="))
		})
	})

	var _ = Describe("Test check password hash", func() {
		It("no valid cases", func() {
			Expect(it.hesh.CheckPasswordHash("123456", "123456")).To(BeFalse())
			Expect(it.hesh.CheckPasswordHash("123456", "")).To(BeFalse())
			Expect(it.hesh.CheckPasswordHash("29l4zNu+i23nf2s3td+bW2Kn6JKlAcO1PqoWsIOL1e0=", "123456")).To(BeFalse())
			Expect(it.hesh.CheckPasswordHash("", "")).To(BeFalse())
		})

		It("valid cases", func() {
			Expect(it.hesh.CheckPasswordHash("123456", "29l4zNu+i23nf2s3td+bW2Kn6JKlAcO1PqoWsIOL1e0=")).To(BeTrue())
			Expect(it.hesh.CheckPasswordHash("test", "Q7DO+ZJl+eNMEOqdNQGSbSezn1fG1nRWHYuiNueoGfs=")).To(BeTrue())
		})
	})

	var _ = Describe("Test jsonify", func() {
		testRows, err := store.Test().GetTestRows();
		if err != nil {
			log.Fatal(err)
		}

		AfterEach(func() {
			defer testRows.Close()
		})

		result := helpers.Jsonify(testRows)

		It("valid cases", func() {
			Expect(result[0]).NotTo(Equal(``))
		})

		It("valid cases", func() {
			Expect(result).NotTo(BeNil())
			Expect(result[0]).To(Equal(`{"id":0,"name":"test_1"}`))
			Expect(result[1]).To(Equal(`,`))
			Expect(result[2]).To(Equal(`{"id":1,"name":"test_2"}`))
			Expect(result[3]).To(Equal(`,`))
			Expect(result[4]).To(Equal(`{"id":2,"name":"test_3"}`))
		})

	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation Suite")
}