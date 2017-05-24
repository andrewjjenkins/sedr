package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/juju/persistent-cookiejar"
	"net/http"
	urlpackage "net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultBase    = "https://companion.orerve.net"
	DefaultLogin   = "/user/login"
	DefaultVerify  = "/user/confirm"
	DefaultRoot    = "/"
	DefaultProfile = "/profile"

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

// Most redirections are allowed but we signal a few with special errors
// to make the authentication workflow correct.
func checkRedirect(req *http.Request, via []*http.Request) error {
	/*
		fmt.Printf("asked to redirect: %s ||", req.URL)
		for _, viaReq := range via {
			fmt.Printf(" %s", viaReq.URL)
		}
		fmt.Printf("\n")
		fmt.Printf("Response headers:\n")
		if req.Response != nil {
			for name, val := range req.Response.Header {
				fmt.Printf("  %s: %s\n", name, val)
			}
		} else {
			fmt.Printf("Response is nil\n")
		}
		for i, viaEntry := range via {
			fmt.Printf("Via %d response headers:\n", i)
			if viaEntry.Response != nil {
				for name, val := range viaEntry.Response.Header {
					fmt.Printf("  %s: %s\n", name, val)
				}
			} else {
				fmt.Printf("Response is nil\n")
			}
		}
	*/

	if req.URL.EscapedPath() == DefaultVerify &&
		len(via) == 1 &&
		via[0].URL.EscapedPath() == DefaultLogin {
		// We just tried to login with user/pass.  It was accepted: we were redirected
		// to the token-verify page.
		return http.ErrUseLastResponse
	} else if req.URL.EscapedPath() == DefaultRoot &&
		len(via) == 1 &&
		via[0].URL.EscapedPath() == DefaultVerify {
		// We just tried to verify the token.  It was accepted: we were redirected to /
		return http.ErrUseLastResponse
	}

	return nil
}

func NewEDClient() (EDClient, error) {
	var err error

	client := EDClient{
		Base: DefaultBase,

		Client: http.Client{
			Timeout: time.Second * 10,

			CheckRedirect: checkRedirect,
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

func (c *EDClient) PostForm(path string, formVals *urlpackage.Values) (*http.Response,
	error) {
	url := c.Base + path

	req, err := http.NewRequest("POST", url, strings.NewReader(formVals.Encode()))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", DefaultUA)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	callHttpDebugger(c, "Posting to "+path, *res)
	return res, nil
}

func (c *EDClient) Get(path string) (*http.Response, error) {
	url := c.Base + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", DefaultUA)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	callHttpDebugger(c, "Get "+path, *res)
	return res, nil
}

func (c *EDClient) GetJSON(path string, out interface{}) error {
	res, err := c.Get(path)
	if err != nil {
		return err
	}

	contentType := res.Header.Get("content-type")
	if contentType != "application/json" {
		return errors.New("GetJSON for " + path + " got content-type " +
			contentType)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

func (c *EDClient) tweakCompanionAppCookie(res *http.Response) {
	// The returned cookie has no persistence or maxage parameter.  Hack one.
	// Also set secure to true.
	for _, cookie := range res.Cookies() {
		if cookie.Name == "CompanionApp" {
			if cookie.MaxAge == 0 && cookie.Expires.IsZero() {
				fmt.Printf("cookie has no expiration, setting one")

				// I have no idea if the cookie is actually valid this long, but I know that
				// I rarely have to reauthenticate any of my clients.
				cookie.Expires = time.Now().Add(365 * 24 * 3600 * time.Second)
			}
			cookie.Secure = true
			c.Jar.SetCookies(res.Request.URL, []*http.Cookie{cookie})
		}
	}
}

func (c *EDClient) Login(email string, password string) error {
	bodyData := urlpackage.Values{}
	bodyData.Set("email", email)
	bodyData.Set("password", password)

	res, err := c.PostForm(DefaultLogin, &bodyData)
	if err != nil {
		return err
	}
	c.tweakCompanionAppCookie(res)
	return nil
}

func (c *EDClient) NeedLogin() bool {
	url := c.Base + DefaultLogin

	parsedUrl, err := urlpackage.Parse(url)
	if err != nil {
		return true
	}
	if len(c.Jar.Cookies(parsedUrl)) != 0 {
		return false
	} else {
		return true
	}
}

func (c *EDClient) Verify(code string) error {
	verifyData := urlpackage.Values{}
	verifyData.Set("code", code)

	res, err := c.PostForm(DefaultVerify, &verifyData)
	if err != nil {
		return err
	}
	c.tweakCompanionAppCookie(res)
	return nil
}

func (c *EDClient) VerifyKeyboard() error {
	var code string
	fmt.Print("Enter verification code from email: ")
	fmt.Scanln(&code)

	return c.Verify(code)
}

func (c *EDClient) SaveJar() error {
	return c.Jar.Save()
}
