package main

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	"github.com/labstack/echo/engine/standard"
	. "github.com/onsi/gomega"
	"github.com/parnurzeal/gorequest"
)

var testURL = "http://localhost:8082"

func Test(t *testing.T) {
	// Init Goblin
	g := Goblin(t)

	// Init Gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	// Init the Web Server
	initDatabaseConnection()
	router := initRouter()
	go router.Run(standard.New(":8082"))

	g.Describe("Integration Tests", func() {
		g.Describe("POST /shorten", func() {
			g.It("should fail to create a shortcode without POSTing any data", func(done Done) {
				request := gorequest.New()
				resp, body, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(Equal(`{"message":"Bad Request"}`))
        done()
			})
			g.It("should fail to create a shortcode without POSTing a URL param", func(done Done) {
				request := gorequest.New()
				resp, body, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(`{"adrian":"impraise"}`).
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(body).To(Equal(`{"message":"url is not present"}`))
        done()
			})
			g.It("should fail to create a shortcode if the desired one does not respect the format", func(done Done) {
				request := gorequest.New()
				resp, body, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(`{"url":"https://impraise.com", "shortcode":"a!@#$%"}`).
					End()
				Expect(resp.StatusCode).To(Equal(422))
				Expect(body).To(Equal(`{"message":"The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{4,}$"}`))
        done()
			})
			g.It("should create a new shortened URL", func(done Done) {
				request := gorequest.New()
				result := new(shorty)
				resp, _, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(`{"url":"https://impraise.com"}`).
					EndStruct(&result)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(result.ShortCode).To(Not(BeNil()))
				Expect(len(*result.ShortCode)).To(Equal(6))
        done()
			})
			g.It("should create a new shortened URL with the prefered shortcode", func(done Done) {
				request := gorequest.New()
				result := new(shorty)
				shortcode := generateRandomString(6)

				resp, _, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(fmt.Sprintf(`{"url":"https://impraise.com", "shortcode":"%s"}`, shortcode)).
					EndStruct(&result)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(result.ShortCode).To(Not(BeNil()))
				Expect(len(*result.ShortCode)).To(Equal(6))
				Expect(*result.ShortCode).To(Equal(shortcode))
        done()
			})
			g.It("should fail to register the same shortcode twice", func(done Done) {
				request := gorequest.New()
				result := new(shorty)
				shortcode := generateRandomString(6)

				resp, _, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(fmt.Sprintf(`{"url":"https://impraise.com", "shortcode":"%s"}`, shortcode)).
					EndStruct(&result)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				Expect(result.ShortCode).To(Not(BeNil()))
				Expect(len(*result.ShortCode)).To(Equal(6))
				Expect(*result.ShortCode).To(Equal(shortcode))

				// Try to register again the same shortcode
				request = gorequest.New()
				resp, body, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(fmt.Sprintf(`{"url":"https://impraise.com", "shortcode":"%s"}`, shortcode)).
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusConflict))
				Expect(body).To(Equal(`{"message":"The the desired shortcode is already in use"}`))
        done()
			})
		})
		g.Describe("GET /:shortcode", func() {
			g.It("should fail to redirect to an unknown shortcode", func(done Done) {
				request := gorequest.New()
				resp, _, _ := request.
					Get(testURL+"/what").
					Set("Content-Type", "application/json").
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
        done()
			})
			g.It("should redirect sucessfully to an url", func(done Done) {
				// Register a new shortcode
				newShortCodeRequest := gorequest.New()
				result := new(shorty)
				shortcode := generateRandomString(6)

				resp, _, _ := newShortCodeRequest.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(fmt.Sprintf(`{"url":"https://impraise.com", "shortcode":"%s"}`, shortcode)).
					EndStruct(&result)

				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				// Now redirect to it

				getShortCodeRequest := gorequest.New()
				resp, _, _ = getShortCodeRequest.
					Get(fmt.Sprintf("%s/%s", testURL, shortcode)).
					Set("Content-Type", "application/json").
					End()

				// // TODO: resp.StatusCode should be http.StatusFound
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
        done()
			})
		})
		g.Describe("GET /:shortcode/stats", func() {
			g.It("should fail to get the stats of an unknown shortcode", func(done Done) {
				request := gorequest.New()
				shortcode := generateRandomString(6)
				resp, _, _ := request.
					Get(fmt.Sprintf("%s/%s/stats", testURL, shortcode)).
					Set("Content-Type", "application/json").
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
        done()
			})
			g.It("should get the stats of a shortcode who's never been visited", func(done Done) {
				request := gorequest.New()
				shortcode := generateRandomString(6)
				result := new(shorty)
				resp, _, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(fmt.Sprintf(`{"url":"https://impraise.com", "shortcode":"%s"}`, shortcode)).
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				resp, _, _ = request.
					Get(fmt.Sprintf("%s/%s/stats", testURL, shortcode)).
					Set("Content-Type", "application/json").
					EndStruct(&result)

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(*result.RedirectCount).To(BeZero())
				Expect(result.LastSeenDate).To(BeNil())
        done()
			})
			g.It("should get the stats of a shortcode who's been visited", func(done Done) {
				request := gorequest.New()
				shortcode := generateRandomString(6)
				result := new(shorty)

				// Register new shortcode
				resp, _, _ := request.
					Post(testURL+"/shorten").
					Set("Content-Type", "application/json").
					Send(fmt.Sprintf(`{"url":"https://impraise.com", "shortcode":"%s"}`, shortcode)).
					End()
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				// Visit shortcode
				getShortCodeRequest := gorequest.New()
				resp, _, _ = getShortCodeRequest.
					Get(fmt.Sprintf("%s/%s", testURL, shortcode)).
					Set("Content-Type", "application/json").
					End()

				// // TODO: resp.StatusCode should be http.StatusFound
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				// Get the stats of the shortcode
				request = gorequest.New()
				resp, _, _ = request.
					Get(fmt.Sprintf("%s/%s/stats", testURL, shortcode)).
					Set("Content-Type", "application/json").
					EndStruct(&result)

				var expectedVisitedCount uint64 = 1
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(*result.RedirectCount).To(Equal(expectedVisitedCount))
				Expect(result.LastSeenDate).To(Not(BeNil()))
        done()
			})
		})
	})
}
