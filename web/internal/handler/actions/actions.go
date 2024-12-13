package actions

import (
	"context"
	"net/http"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/utility"

	"github.com/gin-gonic/gin"
)

func ActionsListGetHandler(c *gin.Context) {
	ctx := context.Background()
	actions, err := repository.ActionsRepositoryInstance.GetAllActions(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "actions_list.html", utility.MergeMaps(*baseH, gin.H{
		"title":   "Actions",
		"actions": actions}))
}
