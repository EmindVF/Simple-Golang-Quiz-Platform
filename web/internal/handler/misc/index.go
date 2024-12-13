package misc

import (
	"net/http"
	"quiz_platform/internal/utility"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "index.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Voprosnja"}))
}
