package utility

import "github.com/gin-gonic/gin"

func MergeMaps(a, b gin.H) gin.H {
	mergedMap := make(gin.H)
	for k, v := range a {
		mergedMap[k] = v
	}
	for k, v := range b {
		mergedMap[k] = v
	}

	return mergedMap
}
