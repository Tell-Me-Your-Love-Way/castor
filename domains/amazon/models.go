package amazon

import "github.com/gin-gonic/gin"

type GetItemsRequest struct {
	Id           string `json:"id"`
	AssociateTag string `json:"associate_tag"`
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
}

var InvalidRequestResponse = gin.H{"error": "Invalid request"}
var InternalServerErrorResponse = gin.H{"error": "Error Quering PA-API"}
