package delta_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	_ "github.com/lib/pq" // ...
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {

	type response struct {
		Data interface{}
		Msg string
	}

	var _ = Describe("delta server test", func() {
        
    Context("test not auth", func() {
				It("No token", func() {
						resp, err := http.Get("http://localhost:4444/api/v1/user")
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
				})
		})
		
		Context("auth test", func() {

				type resData struct {
					Token string
				}

				It("bed auth data", func() {
						r := strings.NewReader(`{"authData": "Wmh1bGV2"}`)
						resp, err := http.Post("http://localhost:4444/api/v1/auth/login", "application/json", r)

						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				})

        It("user login test", func() {
						r := strings.NewReader(`{"authData": "Wmh1bGV2LUZBOjFoYjY3M3Zh"}`)
						resp, err := http.Post("http://localhost:4444/api/v1/auth/login", "application/json", r)

						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						rd := &resData{}
						data := &response{Data: rd}
						err = json.NewDecoder(resp.Body).Decode(data)

						Expect(err).ToNot(HaveOccurred())
						Expect(data.Msg).To(Equal("Login success"))
						Expect(rd.Token).NotTo(BeNil())
        })
    })

	})


	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}
