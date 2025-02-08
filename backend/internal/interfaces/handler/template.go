package handler

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func TemplateHandler(filename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]interface{}{
			"Host": c.Request.Host,
		}

		if authCookie, err := c.Cookie("auth"); err == nil {
			decoded, err := base64.StdEncoding.DecodeString(authCookie)
			if err == nil {
				var userData map[string]string
				if jsonErr := json.Unmarshal(decoded, &userData); jsonErr == nil {
					data["UserData"] = userData
				}
			}
		}

		c.HTML(200, filename, data)
	}
}
