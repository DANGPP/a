package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// DB giả lập
var store = struct {
	sync.RWMutex
	data map[string]string
}{data: make(map[string]string)}

// Middleware log request
func RequestLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"status":     status,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency_ms": latency.Milliseconds(),
		}).Info("HTTP request completed")
	}
}

func main() {
	// Logger setup
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("Failed to log to file, using default stderr")
	}

	// Gin setup
	r := gin.Default()
	r.Use(RequestLogger(log))

	// POST - tạo key-value
	r.POST("/kv", func(c *gin.Context) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		store.Lock()
		store.data[req.Key] = req.Value
		store.Unlock()

		reqID, _ := c.Get("request_id")
		log.WithFields(logrus.Fields{
			"request_id": reqID,
			"action":     "create",
			"key":        req.Key,
			"value":      req.Value,
		}).Info("Key-Value created")

		c.JSON(http.StatusCreated, gin.H{"message": "created"})
	})

	// GET - xem value
	r.GET("/kv/:key", func(c *gin.Context) {
		key := c.Param("key")
		store.RLock()
		value, exists := store.data[key]
		store.RUnlock()

		reqID, _ := c.Get("request_id")
		if !exists {
			log.WithFields(logrus.Fields{
				"request_id": reqID,
				"action":     "read",
				"key":        key,
			}).Warn("Key not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}

		log.WithFields(logrus.Fields{
			"request_id": reqID,
			"action":     "read",
			"key":        key,
			"value":      value,
		}).Info("Key-Value retrieved")

		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	// PUT - sửa value
	r.PUT("/kv/:key", func(c *gin.Context) {
		key := c.Param("key")
		var req struct {
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		store.Lock()
		_, exists := store.data[key]
		if exists {
			store.data[key] = req.Value
		}
		store.Unlock()

		reqID, _ := c.Get("request_id")
		if !exists {
			log.WithFields(logrus.Fields{
				"request_id": reqID,
				"action":     "update",
				"key":        key,
			}).Warn("Key not found for update")
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}

		log.WithFields(logrus.Fields{
			"request_id": reqID,
			"action":     "update",
			"key":        key,
			"value":      req.Value,
		}).Info("Key-Value updated")

		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	})

	// DELETE - xóa key
	r.DELETE("/kv/:key", func(c *gin.Context) {
		key := c.Param("key")
		store.Lock()
		_, exists := store.data[key]
		if exists {
			delete(store.data, key)
		}
		store.Unlock()

		reqID, _ := c.Get("request_id")
		if !exists {
			log.WithFields(logrus.Fields{
				"request_id": reqID,
				"action":     "delete",
				"key":        key,
			}).Warn("Key not found for delete")
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}

		log.WithFields(logrus.Fields{
			"request_id": reqID,
			"action":     "delete",
			"key":        key,
		}).Info("Key-Value deleted")

		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})

	log.Info("Server starting...")
	r.Run(":8080")
}
