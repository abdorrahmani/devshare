package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func okNext() (http.Handler, *bool) {
	called := false
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	return h, &called
}

func TestAuthMiddleware_RejectsWithoutCookie(t *testing.T) {
	next, called := okNext()
	h := AuthMiddleware(next, "secret")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if *called {
		t.Fatal("next should not be called for an unauthenticated request")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "<form") {
		t.Fatal("login form not rendered")
	}
}

func TestAuthMiddleware_AcceptsCorrectPassword(t *testing.T) {
	next, _ := okNext()
	h := AuthMiddleware(next, "secret")
	form := url.Values{"password": {"secret"}}
	req := httptest.NewRequest(http.MethodPost, "/app", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("want 302 redirect, got %d", rr.Code)
	}
	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("no session cookie set")
	}
	c := cookies[0]
	if c.Value == "secret" {
		t.Fatal("cookie stores the plaintext password; expected an opaque token")
	}
	if len(c.Value) < 16 {
		t.Fatalf("session token looks weak: %q", c.Value)
	}
	if !c.HttpOnly || c.SameSite != http.SameSiteLaxMode {
		t.Fatal("cookie missing HttpOnly/SameSite hardening")
	}
}

func TestAuthMiddleware_ValidSessionPasses(t *testing.T) {
	next, called := okNext()
	h := AuthMiddleware(next, "secret")

	form := url.Values{"password": {"secret"}}
	login := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	login.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	lrr := httptest.NewRecorder()
	h.ServeHTTP(lrr, login)
	cookies := lrr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("no cookie issued on login")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookies[0])
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if !*called || rr.Code != http.StatusOK {
		t.Fatalf("valid session should pass through; called=%v code=%d", *called, rr.Code)
	}
}

func TestAuthMiddleware_RejectsWrongPassword(t *testing.T) {
	next, called := okNext()
	h := AuthMiddleware(next, "secret")
	form := url.Values{"password": {"nope"}}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if *called {
		t.Fatal("next called with a wrong password")
	}
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_NoReflectedXSS(t *testing.T) {
	next, _ := okNext()
	h := AuthMiddleware(next, "secret")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL.Path = `/"><script>alert(1)</script>`
	req.URL.RawQuery = `x="><script>alert(1)</script>`
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	// The page is fully static; if it ever varies with request input, that's a reflection.
	if rr.Body.String() != loginHTML {
		t.Fatal("login page varied with request input; possible reflected XSS")
	}
}
