package magalu

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	RedisClient() *redis.Client
	RenderSite(sku, partnerTag string) (string, string, error)
}
type service struct {
	redisClient *redis.Client
	httpClient  *http.Client
}

func (s *service) HttpClient() *http.Client { return s.httpClient }
func (s *service) RedisClient() *redis.Client {
	return s.redisClient
}
func NewService(redisClientParam *redis.Client) Service {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        400,
			MaxIdleConnsPerHost: 200,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	return &service{
		redisClient: redisClientParam,
		httpClient:  client,
	}
}

// RenderSite fetches and renders a webpage, mimicking human-like browser behavior.
func (s *service) RenderSite(sku, partnerTag string) (string, string, error) {
	// Construct initial URL
	baseURL := fmt.Sprintf("https://www.magazinevoce.com.br/%s/p/%s", partnerTag, sku)

	// Create HTTP request
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

	// Add cookies to simulate a user session
	cookies := []string{
		"MLPARCEIRO=" + partnerTag,
		"noe_freight=AUTO",
		"noe_hub_shipping_enabled=1",
		"toggle_wishlist=false",
		"show_seller_score_above=5",
		"FCCDCF=1",
		"ml2_redirect_8020=0",
		"FCNEC=1",
		"mixer_shipping=AUTO",
		"mixer_hub_shipping=true",
		"toggle_pdp_seller_score=true",
		"toggle_vwo=true",
		"toggle_agatha=true",
		"toggle_ads=true",
		"toggle_new_service_page=true",
		"toggle_quick_click=false",
		"enable_fallback_banner=0",
		"pdp_desk_b2c_mixer=1",
		fmt.Sprintf("__uzma=%s", generateRandomUUID()),
		fmt.Sprintf("__uzmb=%d", time.Now().Unix()),
		"__uzme=5845",
		fmt.Sprintf("__uzmc=%d", randomNumber(100000000000, 999999999999)),
		fmt.Sprintf("__uzmd=%d", time.Now().Unix()),
	}
	req.Header.Set("Cookie", strings.Join(cookies, "; "))

	// Execute request
	resp, err := s.HttpClient().Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf(err) // Log ou trate o erro de fechamento do corpo, se necessário
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

	price := doc.Find(`[data-testid="price-value"]`).First().Text()
	trimmed := strings.TrimPrefix(url, "https://www.magazinevoce.com.br/"+partnerTag+"/")

	// Construct final URL
	finalURL := fmt.Sprintf("https://www.magazineluiza.com.br/%s?partner_id=%s", trimmed, partnerTag)

	return price, finalURL, nil
}

// Helper functions to mimic human-like behavior
func generateRandomUUID() string {
	// Simplified UUID generator for demonstration
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(randomNumber(0, 255))
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func randomNumber(min, max int) int {
	return min + int(time.Now().UnixNano()%int64(max-min+1))
}

var ServiceInstance Service
