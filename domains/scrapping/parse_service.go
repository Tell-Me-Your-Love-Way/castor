package scrapping

import (
	"fmt"
	"strings"
)

func (s *service) ParseUrl(url string) (string, string, string, error) {
	domain := s.Domain().FindStringSubmatch(url)

	if strings.Contains(domain[1], "mercadolivre") {
		if strings.Contains(domain[1], "produto") {
			r := s.Id()["mercadolivre"]
			id := r.FindStringSubmatch(url)
			link := fmt.Sprintf("https://%s/MLB-%s", domain, id[1])
			id[1] = "MLB" + id[1]
			return link, domain[1], id[1], nil
		}

		r := s.Id()["magalu"]
		id := r.FindStringSubmatch(url)
		link := fmt.Sprintf("https://%s/p/%s", domain, id[1])
		return link, domain[1], id[1], nil

	}
	if strings.Contains(domain[1], "magalu") {
		r := s.Id()["magalu"]
		id := r.FindStringSubmatch(url)
		link := fmt.Sprintf("https://%s/p/%s", domain, id[1])
		return link, domain[1], id[1], nil
	}
	err := fmt.Errorf("Not Valid URL: %s", url)
	return "", "", "", err
}
