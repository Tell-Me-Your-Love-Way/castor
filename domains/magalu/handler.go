package magalu

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	req := Request{}
	if c.BindJSON(&req) != nil {
		c.JSON(400, InvalidRequestResponse)
		return
	}
	if req.Sku == "" || req.PartnerTag == "" {
		c.JSON(400, InvalidRequestResponse)
		return
	}
	text, url, err := ServiceInstance.RenderSite(req.Sku, req.PartnerTag)
	if err != nil {
		fmt.Printf("Error rendering site: %v\n", err)
		c.JSON(500, InternalServerErrorResponse)
		return
	}
	if text == "" {
		c.JSON(400, gin.H{"error": "Render Error"})
	}
	c.JSON(200, gin.H{"price": text, "url": url})
}
