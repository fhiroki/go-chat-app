package handler

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploaderHandler(c *gin.Context) {
	userID := c.PostForm("user_id")
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	filename := filepath.Join("assets", "avatars", userID+filepath.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Successful"})
}
