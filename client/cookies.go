package client

import (
	"github.com/juju/persistent-cookiejar"
)

const (
	DefaultCookieJarFilename = ".sedr-cookies"
)

func OpenCookieJar(filename string) (*cookiejar.Jar, error) {
	if len(filename) == 0 {
		filename = DefaultCookieJarFilename
	}

	jarOptions := cookiejar.Options{
		Filename: filename,
	}

	return cookiejar.New(&jarOptions)
}
