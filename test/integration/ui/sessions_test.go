package ui

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_client"
	"github.com/supergiant/supergiant/test/fake_core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSessionsList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI Sessions List works correctly", t, func() {

		table := []struct {
			// Mocks
			mockAuthenticated bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				mockAuthenticated:  true,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions",
			},
			// Unauthenticated
			{
				mockAuthenticated:  false,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/new",
			},
		}

		for _, item := range table {

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					if item.mockAuthenticated {
						return &client.Client{}
					}
					return nil
				},
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/sessions", nil)

			// As long as we have a cookie with the right name, it will trigger the
			// use of our fake_core.Sessions above.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			resp, _ := http.DefaultClient.Do(req)

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
		}
	})
}

//------------------------------------------------------------------------------

func TestSessionsNew(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI Sessions New works correctly", t, func() {

		table := []struct {
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/new",
			},
		}

		for _, item := range table {

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/sessions/new", nil)

			resp, _ := http.DefaultClient.Do(req)

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
		}
	})
}

//------------------------------------------------------------------------------

func TestSessionsCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI Sessions Create works correctly", t, func() {

		table := []struct {
			// Input
			inputUsername string
			inputPassword string
			// Mocks
			mockClientCreateError error
			// Expectations
			redirected          bool // how to distinguish between a successful create and a failure
			modelPassedToCreate *model.Session
			responseStatusCode  int
			responseURL         string
		}{
			// A successful example
			{
				inputUsername:         "username",
				inputPassword:         "password",
				mockClientCreateError: nil,
				redirected:            true,
				modelPassedToCreate: &model.Session{
					User: &model.User{
						Username: "username",
						Password: "password",
					},
				},
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions",
			},
			// Client Create error
			{
				inputUsername:         "username",
				inputPassword:         "password",
				mockClientCreateError: errors.New("something bad"),
				modelPassedToCreate: &model.Session{
					User: &model.User{
						Username: "username",
						Password: "password",
					},
				},
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions",
			},
		}

		for _, item := range table {

			var modelPassedToCreate *model.Session
			var redirected bool

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return &client.Client{
					Sessions: &fake_client.Sessions{
						Collection: fake_client.Collection{
							CreateFn: func(m model.Model) error {
								modelPassedToCreate = m.(*model.Session)
								return item.mockClientCreateError
							},
						},
					},
				}
			}

			// This is the authenticated Client for the post-create redirect
			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					return &client.Client{}
				},
			}

			form := url.Values{}
			form.Add("username", item.inputUsername)
			form.Add("password", item.inputPassword)
			body := strings.NewReader(form.Encode())

			req, _ := http.NewRequest("POST", "http://localhost:10000/ui/sessions", body)

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// NOTE we have to mock the setting of the cookie by the Sessions
			// controller here.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			client := http.DefaultClient

			// https://github.com/golang/go/issues/4800
			client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return errors.New("too many redirects")
				}
				if len(via) == 0 {
					return nil
				}

				redirected = true

				for attr, val := range via[0].Header {
					if _, ok := r.Header[attr]; !ok {
						r.Header[attr] = val
					}
				}
				return nil
			}

			resp, _ := client.Do(req)

			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(b))

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
			So(redirected, ShouldEqual, item.redirected)
			So(modelPassedToCreate, ShouldResemble, item.modelPassedToCreate)
		}
	})
}

//------------------------------------------------------------------------------

func TestSessionsGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI Sessions Get works correctly", t, func() {

		table := []struct {
			// Mocks
			mock404 bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/1",
			},
			// 404
			{
				mock404:            true,
				responseStatusCode: 500, // TODO ------------ we probably want a 404 here like the other models
				responseURL:        "http://localhost:10000/ui/sessions/1",
			},
		}

		for _, item := range table {

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return &client.Client{
					Sessions: &fake_client.Sessions{
						Collection: fake_client.Collection{
							GetFn: func(id interface{}, m model.Model) error {
								if item.mock404 {
									return errors.New("404")
								}
								return nil
							},
						},
					},
				}
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/sessions/1", nil)

			resp, _ := http.DefaultClient.Do(req)

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
		}
	})
}
