package common

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type CookieJar struct {
	Data map[string]*http.Cookie
}

func (c *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, val := range cookies {
		if val.Value == "" {
			delete(c.Data, val.Name)
		} else {
			c.Data[val.Name] = val
		}
	}
}

func (c *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	cookies := make([]*http.Cookie, 0, len(c.Data))
	for _, val := range c.Data {
		cookies = append(cookies, val)
	}
	return cookies
}

func (c *CookieJar) Ouput() ([]byte, error) {
	return json.Marshal(c.Data)
}
func (c *CookieJar) Input(data []byte) error {
	return json.Unmarshal(data, &c.Data)
}
func (c *CookieJar) String() string {
	res := ""
	for _, ck := range c.Data {
		res += ck.Name + "=" + ck.Value + ";"
	}
	return res
}
func NewCookie() *CookieJar {
	return &CookieJar{
		Data: make(map[string]*http.Cookie),
	}
}
