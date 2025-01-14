package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql" // Replace with your DB driver (e.g., postgres, sqlite, etc.)
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// Initialize the database connection (replace with your connection string)
	db, err := gorm.Open(mysql.Open("admin:admin@tcp(127.0.0.1:5432)/mydb"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Output is still required, use os.DevNull if needed
			logger.Config{
				LogLevel: logger.Silent, // Disable all logging
			},
		),
	})
	if err != nil {
		panic(err)
	}

	// AutoMigrate the schema
	err = db.AutoMigrate(&Role{}, &Permission{}, &Object{}, &RolePermission{})
	if err != nil {
		panic(err)
	}

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard // Optional: suppress error logs as well
	r := gin.New()

	r.GET("/checkPermission", func(c *gin.Context) {
		type Req struct {
			RoleName       string
			PermissionName string
			ObjectName     string
		}

		var req Req
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var count int64
		err := db.Model(&RolePermission{}).
			Joins("JOIN roles ON roles.ID = role_permissions.RoleID").
			Joins("JOIN permissions ON permissions.ID = role_permissions.PermissionID").
			Joins("JOIN objects ON objects.ID = role_permissions.ObjectID").
			Where("roles.Name = ? AND permissions.Name = ? AND objects.Name = ?", req.RoleName, req.PermissionName, req.ObjectName).
			Count(&count).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			return
		}

		if count > 0 {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusForbidden)
		}
	})

	r.POST("/addPermission", func(c *gin.Context) {
		type Req struct {
			RoleName       string
			PermissionName string
			ObjectName     string
		}

		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find or create the role
		var role Role
		if err := db.FirstOrCreate(&role, Role{Name: req.RoleName}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create role"})
			return
		}

		// Find or create the permission
		var permission Permission
		if err := db.FirstOrCreate(&permission, Permission{Name: req.PermissionName}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create permission"})
			return
		}

		// Find or create the object
		var object Object
		if err := db.FirstOrCreate(&object, Object{Name: req.ObjectName}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find or create object"})
			return
		}

		// Create the role-permission-object association
		rolePermission := RolePermission{
			RoleID:       role.ID,
			PermissionID: permission.ID,
			ObjectID:     object.ID,
		}
		if err := db.Create(&rolePermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role-permission-object association"})
			return
		}

		c.Status(http.StatusOK)
	})

	fmt.Println("run...")
	r.Run(":8080")
}
