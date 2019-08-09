package node

import (
	"github.com/gin-gonic/gin"
	"github.com/nknorg/nkn/node"
	"net/http"
)

func NeighborRouter(router *gin.RouterGroup) {
	router.GET("/node/neighbors", func(context *gin.Context) {
		localNode, exists := context.Get("localNode")

		if exists {
			list := localNode.(*node.LocalNode).GetNeighborInfo()
			context.JSON(http.StatusOK, list)
			return
		}

	})
}
