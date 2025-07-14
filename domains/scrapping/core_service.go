package scrapping

import "regexp"

type Service interface {
	ParseUrl(string) (string, string, string, error)
}
type service struct {
	domain *regexp.Regexp
	id     map[string](*regexp.Regexp)
}

func (s *service) Domain() *regexp.Regexp {
	return s.domain
}
func (s *service) Id() map[string]*regexp.Regexp {
	return s.id
}

func NewService() Service {
	productId := make(map[string]*regexp.Regexp)
	productId["magalu"] = regexp.MustCompile(`/p/([A-Za-z0-9]+)`)
	productId["mercadolivre"] = regexp.MustCompile(`/MLB-(\d+)-`)
	return &service{
		domain: regexp.MustCompile(`https?://([^/]+)`),
		id:     productId,
	}
}

var ServiceInstance Service
