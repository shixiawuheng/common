package lanzouw

import (
	"net/http"
	"net/url"
)

type CookieJar struct {
	cookies map[string]*http.Cookie
}

func (v *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, val := range cookies {
		v.cookies[val.Name] = val
	}
}

func (v *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	cookies := make([]*http.Cookie, 0, len(v.cookies))
	for _, val := range v.cookies {
		cookies = append(cookies, val)
	}
	return cookies
}

func newCookieJar() *CookieJar {
	return &CookieJar{
		cookies: make(map[string]*http.Cookie),
	}
}
