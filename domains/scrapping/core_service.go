package scrapping

import (
	"net/http"
	"regexp"
	"time"
)

type Service interface {
	DomainRegex() *regexp.Regexp
	ItemIdRegex() map[string]*regexp.Regexp
	ClientHttp() *http.Client
	CategoryRegex() map[string]*regexp.Regexp
	ParseUrl(string) (string, string, string, error)
	RenderSite(string, string, string) (string, string, error)
}
type service struct {
	domainRegex   *regexp.Regexp
	itemIdRegex   map[string](*regexp.Regexp)
	categoryRegex map[string]*regexp.Regexp
	clientHttp    *http.Client
}

func (s *service) DomainRegex() *regexp.Regexp {
	return s.domainRegex
}
func (s *service) ItemIdRegex() map[string]*regexp.Regexp {
	return s.itemIdRegex
}
func (s *service) CategoryRegex() map[string]*regexp.Regexp {
	return s.categoryRegex
}
func (s *service) ClientHttp() *http.Client {
	return s.clientHttp
}

func NewService() Service {
	// Initialize regex patterns for item IDs
	itemId := make(map[string]*regexp.Regexp)
	itemId["magalu"] = regexp.MustCompile(`/p/([A-Za-z0-9]+)`)
	itemId["mercadolivre"] = regexp.MustCompile(`/MLB-(\d+)-`)
	category := make(map[string]*regexp.Regexp)
	category["magalu"] = regexp.MustCompile(`https?://[^/]+/[^/]+/p/[^/]+/([^/]+)/([^/]+)`)
	// Initialize HTTP client with custom transport settings
	clientWeb := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        400,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	// Return a new service instance with the compiled regex and HTTP client
	return &service{
		domainRegex:   regexp.MustCompile(`https?://([^/]+)`),
		itemIdRegex:   itemId,
		categoryRegex: category,
		clientHttp:    clientWeb,
	}
}

var ServiceInstance Service
