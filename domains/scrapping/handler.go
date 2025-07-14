package scrapping

import (
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	link, domain, id, err := ServiceInstance.ParseUrl(request.Url)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"link":   link,
		"domain": domain,
		"id":     id,
	})
}
