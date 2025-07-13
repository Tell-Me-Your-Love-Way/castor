package magalu

import "github.com/gin-gonic/gin"

type Request struct {
	Sku        string `json:"sku"`
	PartnerTag string `json:"partner_tag"`
}

var InvalidRequestResponse = gin.H{"error": "Invalid request"}
var InternalServerErrorResponse = gin.H{"error": "Error Quering PA-API"}
