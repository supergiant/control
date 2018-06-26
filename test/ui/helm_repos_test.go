package ui

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_client"
	"github.com/supergiant/supergiant/test/fake_core"
)

func TestHelmReposList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI HelmRepos List works correctly", t, func() {

		table := []struct {
			// Mocks
			// mockClientListItems []*model.HelmRepo
			// mockClientListError error
			mockAuthenticated bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				mockAuthenticated:  true,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/helm_repos",
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

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/helm_repos", nil)

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

func TestHelmReposNew(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI HelmRepos New works correctly", t, func() {

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
				responseURL:        "http://localhost:10000/ui/helm_repos/new",
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

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/helm_repos/new", nil)

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

func TestHelmReposCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI HelmRepos Create works correctly", t, func() {

		table := []struct {
			// Input
			jsonInput string
			// Mocks
			mockAuthenticated     bool
			mockClientCreateError error
			// Expectations
			redirected          bool // how to distinguish between a successful create and a failure
			modelPassedToCreate *model.HelmRepo
			responseStatusCode  int
			responseURL         string
		}{
			// A successful example
			{
				jsonInput: `{
          "name": "test",
          "url": "www.website.com"
        }`,
				mockAuthenticated:     true,
				mockClientCreateError: nil,
				redirected:            true,
				modelPassedToCreate: &model.HelmRepo{
					Name: "test",
					URL:  "www.website.com",
				},
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/helm_repos",
			},

			// Bad JSON input
			{
				jsonInput: `{
          "name": "te..
        }`,
				mockAuthenticated:     true,
				mockClientCreateError: nil,
				modelPassedToCreate:   nil,
				responseStatusCode:    200,
				responseURL:           "http://localhost:10000/ui/helm_repos",
			},

			// Unauthenticated
			{
				jsonInput: `{
          "name": "test"
        }`,
				mockAuthenticated:     false,
				mockClientCreateError: nil,
				redirected:            true, // to login page
				modelPassedToCreate:   nil,
				responseStatusCode:    200,
				responseURL:           "http://localhost:10000/ui/sessions/new",
			},
			// Client Create error
			{
				jsonInput: `{
          "name": "test"
        }`,
				mockAuthenticated:     true,
				mockClientCreateError: errors.New("something bad"),
				modelPassedToCreate: &model.HelmRepo{
					Name: "test",
				},
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/helm_repos",
			},
		}

		for _, item := range table {

			var modelPassedToCreate *model.HelmRepo
			var redirected bool

			// For unauthenticated Session-based routes
			srv.Core.APIClient = func(authType string, authToken string) *client.Client {
				return new(client.Client)
			}

			srv.Core.Sessions = &fake_core.Sessions{
				ClientFn: func(sessionID string) *client.Client {
					if item.mockAuthenticated {
						return &client.Client{
							HelmRepos: &fake_client.HelmRepos{
								Collection: fake_client.Collection{
									CreateFn: func(m model.Model) error {
										modelPassedToCreate = m.(*model.HelmRepo)
										return item.mockClientCreateError
									},
								},
							},
						}
					}
					return nil
				},
			}

			form := url.Values{}
			form.Add("json_input", item.jsonInput)
			body := strings.NewReader(form.Encode())

			req, _ := http.NewRequest("POST", "http://localhost:10000/ui/helm_repos", body)

			// As long as we have a cookie with the right name, it will trigger the
			// use of our fake_core.Sessions above.
			cookie := &http.Cookie{
				Name:  core.SessionCookieName,
				Value: "fake-session-id",
				Path:  "/",
			}
			req.AddCookie(cookie)

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

			// defer resp.Body.Close()
			// b, _ := ioutil.ReadAll(resp.Body)
			// fmt.Println(string(b))

			So(resp.StatusCode, ShouldEqual, item.responseStatusCode)
			So(resp.Request.URL.String(), ShouldEqual, item.responseURL)
			So(redirected, ShouldEqual, item.redirected)
			So(modelPassedToCreate, ShouldResemble, item.modelPassedToCreate)
		}
	})
}

//------------------------------------------------------------------------------

func TestHelmReposGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("UI HelmRepos Get works correctly", t, func() {

		table := []struct {
			// Mocks
			mockAuthenticated bool
			mock404           bool
			// Expectations
			responseStatusCode int
			responseURL        string
		}{
			// A successful example
			{
				mockAuthenticated:  true,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/helm_repos/1",
			},
			// Unauthenticated
			{
				mockAuthenticated:  false,
				responseStatusCode: 200,
				responseURL:        "http://localhost:10000/ui/sessions/new",
			},
			// 404
			{
				mockAuthenticated:  true,
				mock404:            true,
				responseStatusCode: 404,
				responseURL:        "http://localhost:10000/ui/helm_repos/1",
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
						return &client.Client{
							HelmRepos: &fake_client.HelmRepos{
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
					return nil
				},
			}

			req, _ := http.NewRequest("GET", "http://localhost:10000/ui/helm_repos/1", nil)

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
