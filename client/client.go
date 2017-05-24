package client

import (
	"fmt"
	"github.com/juju/persistent-cookiejar"
	"io/ioutil"
	"net/http"
	urlpackage "net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultBase  = "https://companion.orerve.net"
	DefaultLogin = "/user/login"

	// The user agent must be set to a string that the E:D servers will accept.
	DefaultUA = "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1_2 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Mobile/11D257"
)

type EDClient struct {
	// The beginning of the URL to talk to (useful to override for debugging)
	Base string

	// The HTTP client to use for making requests.
	Client http.Client

	// The persistent cookiejar for holding auth details.  Also present at
	// Client.Jar (but as a http.CookieJar)
	Jar *cookiejar.Jar

	// A function to receive every HTTP request and response for debugging.
	HTTPDebugger func(what string, res http.Response)
}

func NewEDClient() (EDClient, error) {
	var err error

	client := EDClient{
		Base: DefaultBase,

		Client: http.Client{
			Timeout: time.Second * 10,
		},

		HTTPDebugger: dumpHttp,
	}

	base, basePresent := os.LookupEnv("SEDR_BASEURL")
	if basePresent {
		client.Base = base
	}

	client.Jar, err = OpenCookieJar("")
	if err != nil {
		return client, err
	}
	client.Client.Jar = client.Jar

	return client, nil
}

func dumpHttp(what string, res http.Response) {
	req := res.Request
	fmt.Printf("[DEBUG] HTTP %d %s %s\n", res.StatusCode, what, req.URL)

	fmt.Printf("[DEBUG] Request %s %s %s:\n", req.Method, req.URL, req.Proto)
	for name, val := range req.Header {
		fmt.Printf("[DEBUG]     %s: %s\n", name, val)
	}

	fmt.Printf("[DEBUG] Response %s %s:\n", res.Status, res.Proto)
	for name, val := range res.Header {
		fmt.Printf("[DEBUG]     %s: %s\n", name, val)
	}

	fmt.Printf("[DEBUG] ****\n")
}

func callHttpDebugger(c *EDClient, what string, res http.Response) {
	if c.HTTPDebugger != nil {
		c.HTTPDebugger(what, res)
	}
}

func (c *EDClient) Login(email string, password string) error {
	url := c.Base + DefaultLogin

	bodyData := urlpackage.Values{}
	bodyData.Set("email", email)
	bodyData.Set("password", password)
	req, err := http.NewRequest("POST", url, strings.NewReader(bodyData.Encode()))

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", DefaultUA)

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	callHttpDebugger(c, "Login", *res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		fmt.Printf("Login response body: %s\n", body)
	}

	// The returned cookie has no persistence or maxage parameter.  Hack one.
	// Also set secure to true.
	for _, cookie := range res.Cookies() {
		if cookie.Name == "CompanionApp" {
			if cookie.MaxAge == 0 && cookie.Expires.IsZero() {
				fmt.Printf("cookie has no expiration, setting one")
				cookie.Expires = time.Now().Add(3600 * time.Second)
			}
			cookie.Secure = true
			c.Jar.SetCookies(req.URL, []*http.Cookie{cookie})
		}
	}

	if res.StatusCode == 200 {
		err = c.Jar.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
