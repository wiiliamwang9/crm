package config

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupStaticRoutes(r *gin.Engine) {
	webPath := "../web"
	if _, err := os.Stat(webPath); os.IsNotExist(err) {
		fmt.Printf("Warning: Web directory not found at %s\n", webPath)
		webPath = "./web"
	}

	r.Static("/static", webPath)

	// 添加target目录的静态文件服务，自定义处理HTML文件的Content-Type
	targetHandler := func(c *gin.Context) {
		filePath := c.Param("filepath")
		fullPath := filepath.Join(webPath, "target", filePath)

		// 检查文件是否存在
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// 根据文件扩展名设置正确的Content-Type
		ext := filepath.Ext(filePath)
		switch ext {
		case ".html":
			c.Header("Content-Type", "text/html; charset=utf-8")
		case ".css":
			c.Header("Content-Type", "text/css; charset=utf-8")
		case ".js":
			c.Header("Content-Type", "application/javascript; charset=utf-8")
		case ".json":
			c.Header("Content-Type", "application/json; charset=utf-8")
		default:
			c.Header("Content-Type", "application/octet-stream")
		}

		c.File(fullPath)
	}

	r.GET("/target/*filepath", targetHandler)
	r.HEAD("/target/*filepath", targetHandler)

	// 添加config目录的静态文件服务
	configHandler := func(c *gin.Context) {
		filePath := c.Param("filepath")
		fullPath := filepath.Join(webPath, "config", filePath)

		// 检查文件是否存在
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// 根据文件扩展名设置正确的Content-Type
		ext := filepath.Ext(filePath)
		switch ext {
		case ".js":
			c.Header("Content-Type", "application/javascript; charset=utf-8")
		case ".json":
			c.Header("Content-Type", "application/json; charset=utf-8")
		case ".css":
			c.Header("Content-Type", "text/css; charset=utf-8")
		default:
			c.Header("Content-Type", "application/octet-stream")
		}

		c.File(fullPath)
	}

	r.GET("/config/*filepath", configHandler)
	r.HEAD("/config/*filepath", configHandler)

	r.GET("/", func(c *gin.Context) {
		indexPath := webPath + "/index.html"
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Static files not found"})
			return
		}
		c.File(indexPath)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "CRM API",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
		} else {
			indexPath := webPath + "/index.html"
			if _, err := os.Stat(indexPath); os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Static files not found"})
			} else {
				c.File(indexPath)
			}
		}
	})

	fmt.Printf("Web files path: %s\n", webPath)
}
