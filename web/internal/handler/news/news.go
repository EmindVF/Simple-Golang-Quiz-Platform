package news

import (
	"context"
	"fmt"
	"net/http"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/middleware"
	"quiz_platform/internal/utility"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewsListGetHandler(c *gin.Context) {
	ctx := context.Background()
	news, err := repository.NewsRepositoryInstance.GetAllNews(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "news_list.html", utility.MergeMaps(*baseH, gin.H{
		"title": "News",
		"news":  news}))
}

func NewsViewGetHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()
	news, err := repository.NewsRepositoryInstance.GetNewsById(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "news_view.html", utility.MergeMaps(*baseH, gin.H{
		"title": "News",
		"news":  news}))
}

func NewsCreateFormGetHandler(c *gin.Context) {
	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "news_form.html", utility.MergeMaps(*baseH, gin.H{
		"title":  "News",
		"news":   nil,
		"action": "/news"}))
}

func NewsEditFormGetHandler(c *gin.Context) {

	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()
	news, err := repository.NewsRepositoryInstance.GetNewsById(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "news_form.html", utility.MergeMaps(*baseH, gin.H{
		"title":  "News",
		"news":   news,
		"action": fmt.Sprintf("/news/%d", id)}))
}

func NewsDeletePostHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()
	err = repository.NewsRepositoryInstance.DeleteNews(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/news")
}

func NewsCreatePostHandler(c *gin.Context) {
	var userId int32
	data, ok := c.Get("sessionData")
	if ok {
		if sessionData, ok := data.(*middleware.SessionData); ok {
			userId = sessionData.UserId
		}
	}

	title := c.PostForm("title")
	newsText := c.PostForm("news_text")

	ctx := context.Background()
	_, err := repository.NewsRepositoryInstance.AddNews(ctx, title, newsText, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/news")
}

func NewsEditPostHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	title := c.PostForm("title")
	newsText := c.PostForm("news_text")

	ctx := context.Background()
	err = repository.NewsRepositoryInstance.EditNews(ctx, id, title, newsText)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/news")
}
