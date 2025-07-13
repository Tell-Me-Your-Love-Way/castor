package amazon

import "github.com/gin-gonic/gin"

func HandlerQuery(c *gin.Context) {
	req := GetItemsRequest{}
	if c.BindJSON(&req) != nil {
		c.JSON(400, InvalidRequestResponse)
		return
	}
	if req.Id == "" || req.AssociateTag == "" || req.AccessKey == "" || req.SecretKey == "" {
		c.JSON(400, InvalidRequestResponse)
		return
	}
	res, err := ServiceInstance.QueryPAAPI(req.Id, req.AssociateTag, req.AccessKey, req.SecretKey)
	if err != nil {
		c.JSON(500, InternalServerErrorResponse)
		return
	}

	c.JSON(200, gin.H{"response": res})
}
