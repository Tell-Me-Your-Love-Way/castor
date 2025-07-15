package scrapping

import (
	"fmt"
	"regexp"
	"strings"
)

func (s *service) ParseUrl(url string) (string, string, string, error) {
	domainMatch := s.domainRegex.FindStringSubmatch(url)
	if len(domainMatch) != 2 {
		return "", "", "", fmt.Errorf("Domínio não encontrado na URL: %s", url)
	}
	domain := domainMatch[1]
	// --- LOG ---
	fmt.Printf("Domain: %s\n", domain)
	// === LOG ===
	var (
		r  *regexp.Regexp
		id []string
	)
	err := fmt.Errorf("Not Valid URL: %s", url)
	switch {
	case strings.Contains(domain, "mercadolivre.com"):
		if strings.Contains(domain, "produto") {
			r = s.itemIdRegex["mercadolivre"]
			id = r.FindStringSubmatch(url)
			if len(id) != 2 {
				return "", "", "", err
			}
			link := fmt.Sprintf("https://%s/MLB-%s", domain, id[1])
			// --- LOG ---
			fmt.Printf("Url Parsed\n")
			// === LOG ===
			return link, domain, "MLB" + id[1], nil
		}
		r = s.itemIdRegex["magalu"]
		id = r.FindStringSubmatch(url)
		if len(id) != 2 {
			return "", "", "", err
		}
		link := fmt.Sprintf("https://%s/p/%s", domain, id[1])
		// --- LOG ---
		fmt.Printf("Url Parsed\n")
		// === LOG ===
		return link, domain, id[1], nil

	case strings.Contains(domain, "magazineluiza.com"):
		r = s.itemIdRegex["magalu"]
		id = r.FindStringSubmatch(url)
		if len(id) != 2 {
			return "", "", "", err
		}
		cat := s.categoryRegex["magalu"].FindStringSubmatch(url)
		link := fmt.Sprintf("https://%s/produto/p/%s/%s/%s", domain, id[1], cat[1], cat[2])
		// --- LOG ---
		fmt.Printf("Url Parsed\n")
		// === LOG ===
		return link, domain, id[1], nil
	}

	return "", "", "", err
}
