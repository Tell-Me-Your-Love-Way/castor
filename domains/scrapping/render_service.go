package scrapping

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

func (s *service) RenderSite(baseURL, domain, partnerTag string) (string, string, error) {

	// Create HTTP request

	// --- LOG ---
	fmt.Printf("Base URL: %s\n", baseURL)
	// === LOG ===
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	resp, err := s.ClientHttp().Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()
	url := resp.Request.URL.String()
	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Handle response body based on encoding
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer func() {
			if cerr := gz.Close(); cerr != nil {
				fmt.Printf("erro ao fechar gzip reader: %v\n", cerr)
			}
		}()
		reader = gz
	}

	// Read response body
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		return "", "", fmt.Errorf("erro ao criar documento goquery: %w", err)
	}
	var (
		price    string
		finalURL string
	)
	if strings.Contains(domain, "mercadolivre") {
		// For Mercado Livre, we need to set the partner tag in the URL
		price = doc.Find(`[class="andes-money-amount__fraction"]`).First().Text()
		finalURL = fmt.Sprintf("%s?aff_id=%s", url, partnerTag)
	} else if strings.Contains(domain, "magazineluiza") {
		price = doc.Find(`[data-testid="price-value"]`).First().Text()
		finalURL = fmt.Sprintf("%s?partner_id=%s", url, partnerTag)
	}

	// --- LOG ---
	fmt.Printf("URL: %s\n", finalURL)
	// === LOG ===

	return price, finalURL, nil
}
