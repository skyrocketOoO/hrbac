package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Create a Neo4j connection using the Bolt protocol
func createNeo4jDriver() (neo4j.Driver, error) {
	// Replace with your Neo4j connection details
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "admin"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	return driver, nil
}

// API to check if a role has a specific permission for an object
func checkPermission(c *gin.Context) {
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

	// Connect to Neo4j using Bolt protocol
	driver, err := createNeo4jDriver()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer driver.Close()

	// Start a session and run the Cypher query to check permission
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := `
		MATCH (r:Role {name: $role_name})-[:HAS_PERMISSION]->(p:Permission {name: $permission_name})-[:APPLIES_TO]->(o:Object {name: $object_name})
		RETURN r, p, o
	`

	params := map[string]interface{}{
		"role_name":       req.RoleName,
		"permission_name": req.PermissionName,
		"object_name":     req.ObjectName,
	}

	result, err := session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.Next() {
		// If a record is found, the role has the permission for the object
		c.JSON(http.StatusOK, gin.H{"message": "Permission exists"})
		return
	}

	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Err().Error()})
		return
	}

	// If no result found, return a message indicating the permission doesn't exist
	c.JSON(http.StatusNotFound, gin.H{"message": "Permission does not exist"})
}

// API to add a permission for a role and object
func addPermission(c *gin.Context) {
	var req struct {
		RoleName       string
		PermissionName string
		ObjectName     string
	}

	// Parse the request body into the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Connect to Neo4j using Bolt protocol
	driver, err := createNeo4jDriver()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer driver.Close()

	// Start a session and run the Cypher query to add the permission
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	query := `
		MERGE (r:Role {name: $role_name})
		MERGE (p:Permission {name: $permission_name})
		MERGE (o:Object {name: $object_name})
		MERGE (r)-[:HAS_PERMISSION]->(p)-[:APPLIES_TO]->(o)
	`

	params := map[string]interface{}{
		"role_name":       req.RoleName,
		"permission_name": req.PermissionName,
		"object_name":     req.ObjectName,
	}

	_, err = session.Run(query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission added successfully"})
}

func main() {
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard // Optional: suppress error logs as well
	r := gin.New()

	// Route to check if a role has permission for an object
	r.GET("/checkPermission", checkPermission)

	// Route to add a permission for a role and object
	r.POST("/addPermission", addPermission)

	// Start the Gin server
	fmt.Println("run...")
	r.Run(":8080")
}
