package main

import (
	"net/http"
	"net/http/httptest"
	"stndalng/route"
	"testing"

	"gopkg.in/gavv/httpexpect.v2"
)

// JWT token authentication tests.
func loginTestCase(e *httpexpect.Expect) {
	type Login struct {
		Username string `form:"username"`
		Password string `form:"password"`
		Role     string `form:"role"`
	}

	// login with bad username
	e.POST("/api/login").WithForm(Login{"<bad username>", "<bad password>", ""}).
		Expect().
		Status(http.StatusUnauthorized)

	// login with bad password
	e.POST("/api/login").WithForm(Login{"superadmin", "<bad password>", ""}).
		Expect().
		Status(http.StatusUnauthorized)

	// login with bad role
	e.POST("/api/login").WithForm(Login{"superadmin", "pass123", "<bad role>"}).
		Expect().
		Status(http.StatusUnauthorized)

	// login with right username, password and role
	e.POST("/api/login").WithForm(Login{"superadmin", "pass123", "root"}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	// login with right password and role empty
	r := e.POST("/api/login").WithForm(Login{"superadmin", "pass123", ""}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	r.Keys().ContainsOnly("code", "token")

	token := r.Value("token").String().Raw()

	e.GET("/api/test/restricted").
		Expect().
		Status(http.StatusBadRequest)

	e.GET("/api/test/restricted").WithHeader("Authorization", "Bearer <bad token>").
		Expect().
		Status(http.StatusUnauthorized)

	e.GET("/api/test/restricted").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).Body().Equal("hello, world!")

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	auth.GET("/api/test/restricted").
		Expect().
		Status(http.StatusOK).Body().Equal("hello, world!")
}

func TestRoutes(t *testing.T) {
	// Setup
	handler := route.Init()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	loginTestCase(e)
}
