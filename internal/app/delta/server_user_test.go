package delta_test

import (
	"encoding/json"
	"net/http"
	"strconv"
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

	const AUTH_DATA string = "Wmh1bGV2LUZBOjFoYjY3M1Zh"

	var (
		token string
		id int
		newPass string
	)

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
						r := strings.NewReader(`{"authData": "` + AUTH_DATA + `"}`)
						resp, err := http.Post("http://localhost:4444/api/v1/auth/login", "application/json", r)

						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						rd := &resData{}
						data := &response{Data: rd}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("Login success"))
						Expect(rd.Token).NotTo(BeNil())
						token = rd.Token
        })
		})

		Context("Create user", func() {

				type resData struct {
					Login string
					Name string
					JobCode string
					Email string
					Phone string
					ID int
				}

				It("valid case", func() {
						form := strings.NewReader(`
						{
							"name": "test_user",
							"login": "test_user",
							"password": "NWEd8e3gfh834",
							"jobCode": "root",
							"email": "test@test.ru",
							"phone": "+9999999999"
						}
						`)
						resp, err := http.Post("http://localhost:4444/api/v1/user?jwt=" + token, "application/json", form)

						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusCreated))

						rd := &resData{}
						data := &response{Data: rd}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()
						
						Expect(data.Msg).To(Equal("User created"))

						Expect(rd.Login).To(Equal("test_user"))
						Expect(rd.Name).To(Equal("test_user"))
						Expect(rd.JobCode).To(Equal("root"))
						Expect(rd.Email).To(Equal("test@test.ru"))
						Expect(rd.Phone).To(Equal("+9999999999"))
						
						id = rd.ID
				})
		})

		Context("Replace user user", func() {

				type resData struct {
					Name string
				}

				It("valid case", func() {
						form := strings.NewReader(`
						{
							"name": "test_user_replace",
							"login": "test_user",
							"jobCode": "root",
							"email": "test@test.ru",
							"phone": "+9999999999"
						}
						`)

						client := &http.Client{}
						req, err := http.NewRequest("PUT", "http://localhost:4444/api/v1/user?jwt=" + token, form)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						rd := &resData{}
						data := &response{Data: rd}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("User replace"))
						Expect(rd.Name).To(Equal("test_user_replace"))
				})
		})

		Context("Delete user user", func() {

				It("valid case", func() {
						client := &http.Client{}
						req, err := http.NewRequest("DELETE", "http://localhost:4444/api/v1/user/" + strconv.Itoa(id) + "?jwt=" + token, nil)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						data := &response{}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("User removed"))

				})
		})

		Context("Get user", func() {

				It("valid case", func() {
						client := &http.Client{}
						req, err := http.NewRequest("GET", "http://localhost:4444/api/v1/user?jwt=" + token, nil)
						resp, err := client.Do(req)
		
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						data := &response{}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("Session received"))

				})
		})

		Context("Get user list", func() {

				type resData struct {
					Size int
					Page int
					Result []interface{}
				}

				It("valid case", func() {
						client := &http.Client{}
						req, err := http.NewRequest("GET", "http://localhost:4444/api/v1/user/list/2/0?jwt=" + token, nil)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						rs := &resData{}
						data := &response{Data: rs}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("User list received"))
						Expect(rs.Page).To(Equal(0))
						Expect(rs.Size).To(Equal(2))
						Expect(len(rs.Result)).To(Equal(2))
				})
		})

		Context("Change user password", func() {

				It("valid case, change first", func() {
						form := strings.NewReader(`
						{
								"password": "1hb673va",
								"new": "1hb673va1",
								"confirm": "1hb673va1"
						}
						`)
						client := &http.Client{}
						req, err := http.NewRequest("PUT", "http://localhost:4444/api/v1/user/password?jwt=" + token, form)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						data := &response{}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("Password is changed"))
				})

				It("valid case, change last", func() {
						form := strings.NewReader(`
						{
								"password": "1hb673va1",
								"new": "1hb673va",
								"confirm": "1hb673va"
						}
						`)
						client := &http.Client{}
						req, err := http.NewRequest("PUT", "http://localhost:4444/api/v1/user/password?jwt=" + token, form)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						data := &response{}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("Password is changed"))
				})
		})

		Context("Reset user password", func() {

				type resData struct {
					Login string
					NewPassword string
				}

				It("valid case, reset", func() {
						form := strings.NewReader(`
						{
								"login": "Zhulev-FA"
						}
						`)
						client := &http.Client{}
						req, err := http.NewRequest("POST", "http://localhost:4444/api/v1/user/password?jwt=" + token, form)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						rs := &resData{}
						data := &response{Data: rs}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("Password is reset"))
						Expect(rs.Login).To(Equal("Zhulev-FA"))
						Expect(len(rs.NewPassword)).To(Equal(8))
						newPass = rs.NewPassword
				})

				It("valid case, change last", func() {
						form := strings.NewReader(`
						{
								"password": "`+ newPass +`",
								"new": "1hb673va",
								"confirm": "1hb673va"
						}
						`)
						client := &http.Client{}
						req, err := http.NewRequest("PUT", "http://localhost:4444/api/v1/user/password?jwt=" + token, form)
						resp, err := client.Do(req)
						Expect(err).ToNot(HaveOccurred())
						Expect(resp.StatusCode).To(Equal(http.StatusOK))

						data := &response{}
						err = json.NewDecoder(resp.Body).Decode(data)
						resp.Body.Close()

						Expect(data.Msg).To(Equal("Password is changed"))
				})
		})

	})


	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}
