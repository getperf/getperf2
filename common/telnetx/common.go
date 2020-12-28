package telnetx

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func ParseUrl(uri string) (string, error) {
	if !strings.HasPrefix(uri, "telnet://") {
		return "", nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrapf(err, "parse url %s", uri)
	}
	port := "23"
	if u.Port() != "" {
		port = u.Port()
	}
	return u.Hostname() + ":" + port, nil
}

// 改行コードを統一する。
func convNewline(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
}
