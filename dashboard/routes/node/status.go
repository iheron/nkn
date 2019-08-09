package node

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nknorg/nkn/node"
	"github.com/nknorg/nkn/util/config"
	"github.com/nknorg/nkn/util/log"
	"net/http"
)

type NodeRouter struct {
}

func (nodeRouter *NodeRouter) Router(router *gin.RouterGroup) {
	router.GET("/node/status", func(context *gin.Context) {
		var out map[string]interface{}

		localNode, exists := context.Get("localNode")

		if exists {
			buf, err := localNode.(*node.LocalNode).MarshalJSON()
			if err != nil {
				log.WebLog.Error(err)
				context.AbortWithError(http.StatusInternalServerError, err)
			}
			err = json.Unmarshal(buf, &out)
			if err != nil {
				log.WebLog.Error(err)
				context.AbortWithError(http.StatusInternalServerError, err)
			}
			out["beneficiaryAddr"] = config.Parameters.BeneficiaryAddr
			context.JSON(200, out)
			return
		}

		context.JSON(200, gin.H{
			"status": "...",
		})
	})
}
