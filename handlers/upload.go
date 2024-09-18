package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"iMessage/db"
	"iMessage/utils"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filename, err = filepath.Abs("./uploads/" + filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取 URL 参数
	messageTemplate := c.Query("message_template")
	scheduledTime := c.Query("scheduled_time")
	datetime, err := time.Parse("2006-01-02 15:04:05", scheduledTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid scheduled_time: %s", scheduledTime)})
		return
	}
	if messageTemplate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message_template is empty"})
		return
	}

	err = c.SaveUploadedFile(file, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utils.CsvLoop(filename, func(strings []string) error {
		if len(strings) != 2 {
			return fmt.Errorf("length must be 2: name and phone")
		}
		var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(strings[1]) {
			return fmt.Errorf("phone number format error: %+v", strings[1])
		}
		return nil
	})
	if err != nil {
		os.Remove(filename)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.CreateActivity(uuid.New().String(), messageTemplate, datetime, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!"})
}
