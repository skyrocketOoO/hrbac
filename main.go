package main

import (
	"fmt"
	"log"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

func main() {
	// Connection parameters
	host := "127.0.0.1"
	port := 9669
	username := "root"
	password := "nebula"

	// Create a new Nebula pool
	poolConfig := nebula.GetDefaultConf()
	address := nebula.HostAddress{Host: host, Port: port}
	pool, err := nebula.NewConnectionPool([]nebula.HostAddress{address}, poolConfig, nebula.DefaultLogger{})
	if err != nil {
		log.Fatalf("Failed to initialize the connection pool: %v", err)
	}
	defer pool.Close()

	// Open a session
	session, err := pool.GetSession(username, password)
	if err != nil {
		log.Fatalf("Failed to create a session: %v", err)
	}
	defer session.Release()

	// Execute a query to select a space
	useSpaceQuery := `
		CREATE SPACE IF NOT EXISTS my_space (vid_type=FIXED_STRING(30));
		USE my_space;
	` // Replace with the actual space name
	_, err = session.Execute(useSpaceQuery)
	if err != nil {
		log.Fatalf("Failed to select the space: %v", err)
	}
	fmt.Println("Space selected successfully!")

	// Create schema (tags and edges)
	createSchema := `
    CREATE TAG IF NOT EXISTS object(type string);
    CREATE TAG IF NOT EXISTS subject(type string);
    CREATE EDGE IF NOT EXISTS relation(type string);
	`

	if _, err := session.Execute(createSchema); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}
	fmt.Println("Schema created successfully!")

	// Insert data (objects, subjects, and relations)
	insertData := `
    -- Insert objects
    INSERT VERTEX object(type) VALUES 
        "doc1":("document"), 
        "doc2":("document");

    -- Insert subjects
    INSERT VERTEX subject(type) VALUES 
        "user1":("user"), 
        "user2":("user"), 
        "group1":("group");

    -- Insert relations
    INSERT EDGE relation(type) VALUES 
        "doc1"->"user1":("viewer"), 
        "doc1"->"user2":("editor"), 
        "doc2"->"group1":("owner");
	`

	if _, err := session.Execute(insertData); err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
	fmt.Println("Data inserted successfully!")

	// Query viewers for a specific object
	queryViewers := `
		USE my_space;
    MATCH (o:object)-[r:relation {type: "viewer"}]->(s:subject)
    WHERE o.id == "doc1"
    RETURN s;
	`

	resp, err := session.Execute(queryViewers)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	if !resp.IsSucceed() {
		log.Fatalf("Query failed: %s", resp.GetErrorMsg())
	}
	fmt.Println("Viewers for doc1:")
	for _, row := range resp.GetRows() {
		fmt.Println(row.GetValues())
	}

	// Query documents user1 has access to
	queryDocs := `
    MATCH (s:subject)-[r:relation]->(o:object)
    WHERE s.id == "user1"
    RETURN o, r.type;
	`

	resp, err = session.Execute(queryDocs)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	fmt.Println("Documents user1 has access to:")
	for _, row := range resp.GetRows() {
		fmt.Println(row.GetValues())
	}

	// Add a new viewer
	addViewer := `
    INSERT EDGE relation(type) VALUES "doc2"->"user1":("viewer");
	`

	if _, err := session.Execute(addViewer); err != nil {
		log.Fatalf("Failed to add viewer: %v", err)
	}
	fmt.Println("Viewer added successfully!")

	// Update a relation
	updateRelation := `
    UPDATE EDGE "doc1"->"user2" SET type = "owner";
	`

	if _, err := session.Execute(updateRelation); err != nil {
		log.Fatalf("Failed to update relation: %v", err)
	}
	fmt.Println("Relation updated successfully!")

	// Delete a viewer
	deleteViewer := `
    DELETE EDGE relation "doc1"->"user1";
	`

	if _, err := session.Execute(deleteViewer); err != nil {
		log.Fatalf("Failed to delete viewer: %v", err)
	}
	fmt.Println("Viewer deleted successfully!")

	// Delete an object (vertex)
	deleteObject := `
    DELETE VERTEX "doc1";
	`

	if _, err := session.Execute(deleteObject); err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}
	fmt.Println("Object deleted successfully!")
}
