package amazon

import (
	"context"
	"fmt"
	paapi5 "github.com/goark/pa-api"
	"github.com/goark/pa-api/entity"
	"github.com/goark/pa-api/query"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Service interface {
	RedisClient() *redis.Client
	QueryPAAPI(id, tag, access, secret string) (string, error)
}

type service struct {
	redisClient *redis.Client
}

func (s *service) RedisClient() *redis.Client {
	return s.redisClient
}
func NewService(redisClientParam *redis.Client) Service {
	return &service{
		redisClient: redisClientParam,
	}
}

func (s *service) QueryPAAPI(id, tag, access, secret string) (string, error) {
	client := paapi5.New(
		paapi5.WithMarketplace(paapi5.LocaleBrazil),
	).CreateClient(
		tag,
		access,
		secret,
		paapi5.WithHttpClient(http.DefaultClient),
	)
	q := query.NewGetItems(
		client.Marketplace(),
		client.PartnerTag(),
		client.PartnerType(),
	).ASINs([]string{id}).EnableImages().EnableItemInfo().EnableParentASIN()

	body, err := client.RequestContext(context.Background(), q)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}
	//io.Copy(os.Stdout, bytes.NewReader(body))

	//Decode JSON
	res, err := entity.DecodeResponse(body)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	}
	return res.String(), nil
}

var ServiceInstance Service
