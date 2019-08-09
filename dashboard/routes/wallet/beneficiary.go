package wallet

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	. "github.com/nknorg/nkn/common"
	"github.com/nknorg/nkn/dashboard/auth"
	"github.com/nknorg/nkn/util/config"
	"github.com/nknorg/nkn/util/log"
	"io/ioutil"
	"net/http"
	"os"
)

type SetBeneficiaryData struct {
	BeneficiaryAddr string `form:"beneficiaryAddr" binding:"required"`
}

func BeneficiaryRouter(router *gin.RouterGroup) {
	router.PUT("/current-wallet/beneficiary", auth.WalletAuth(), func(context *gin.Context) {
		var data SetBeneficiaryData
		if err := context.ShouldBind(&data); err != nil {
			log.WebLog.Error(err)
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		_, err := ToScriptHash(data.BeneficiaryAddr)
		if err != nil {
			log.WebLog.Errorf("parse BeneficiaryAddr error: %v", err)
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		configFile := config.ConfigFile
		if configFile == "" {
			configFile = config.DefaultConfigFile
		}
		if _, err := os.Stat(configFile); err == nil {
			file, err := ioutil.ReadFile(configFile)
			if err != nil {
				log.WebLog.Error("Read config file error:", err)
				context.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// Remove the UTF-8 Byte Order Mark
			file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

			var configuration map[string]interface{}
			err = json.Unmarshal(file, &configuration)
			if err != nil {
				log.WebLog.Error(err)
				context.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// set beneficiary address
			configuration["BeneficiaryAddr"] = data.BeneficiaryAddr

			bytes, err := json.MarshalIndent(&configuration, "", "    ")
			if err != nil {
				log.WebLog.Error(err)
				context.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			err = ioutil.WriteFile(configFile, bytes, 0666)
			if err != nil {
				log.WebLog.Error(err)
				context.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			context.JSON(http.StatusOK, gin.H{
				"message":      configuration["BeneficiaryAddr"],
				"currentValue": config.Parameters.BeneficiaryAddr,
			})
			return
		} else {
			log.WebLog.Error("Config file not exists.")
			context.AbortWithError(http.StatusInternalServerError, errors.New("Config file not exists."))
			return
		}

	})

}
