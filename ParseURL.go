package gos

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

//URL embeds net/url and adds extra fields ontop
type URL struct {
	*url.URL
	TLD, ETLDPlus1, DomainName, BaseDomainName, SubDomainName, Port string
}

// parseURL parses a string as a URL and returns a *url.URL
// or any error that occured. If the initially parsed URL
// has no scheme, 'http://' is prepended and the string is
// re-parsed
func parseURL(rawURL string) (parsedURL *url.URL, err error) {
	parsedURL, err = url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme == "" {
		return url.Parse("http://" + rawURL)
	}

	return parsedURL, nil
}

//ParseURL mirrors net/url.Parse except instead it returns
//a tld.URL, which contains extra fields.
func ParseURL(rawURL string) (*URL, error) {
	url, err := parseURL(rawURL)
	if err != nil {
		return nil, err
	}

	if url.Host == "" {
		return &URL{URL: url}, nil
	}

	var domainName, port string

	for i := len(url.Host) - 1; i >= 0; i-- {
		if url.Host[i] == ':' {
			domainName = url.Host[:i]
			port = url.Host[i+1:]
			break
		} else if url.Host[i] < '0' || url.Host[i] > '9' {
			domainName = url.Host
		}
	}

	etldPlus1, err := publicsuffix.EffectiveTLDPlusOne(domainName)
	if err != nil {
		return nil, err
	}

	i := strings.Index(etldPlus1, ".")
	baseDomainName := etldPlus1[0:i]
	tld := etldPlus1[i+1:]

	subDomainName := ""
	if rest := strings.TrimSuffix(domainName, "."+etldPlus1); rest != domainName {
		subDomainName = rest
	}

	return &URL{
		TLD:            tld,
		ETLDPlus1:      etldPlus1,
		DomainName:     domainName,
		SubDomainName:  subDomainName,
		BaseDomainName: baseDomainName,
		Port:           port,
		URL:            url,
	}, nil
}
