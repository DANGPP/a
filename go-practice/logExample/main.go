package main

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"os"
)
var log *logrus.Logger
func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2001-12-31 15:04:05",
	})

	log.SetLevel(logrus.InfoLevel)

	log.SetReportCaller(true)

	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)
}

func main(){
	r:=gin.Default()

	r.Use(func(c *gin.Context) {
		start := time.Now()

		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path": c.Request.URL.Path,
			"query": c.Request.URL.RawQuery,
			"ip" : c.ClientIP(),
		}).Info("Incoming request")
		
	
		c.Next()
		log.WithFields(logrus.Fields{
			"status": c.Writer.Status(),
			"duration": time.Since(start).String(),
		}).Info("response sent")

	})

	r.GET("/ping", func(c *gin.Context) {
		log.WithFields(logrus.Fields{
			"handler": "ping",
		}).Info("Processing request")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}