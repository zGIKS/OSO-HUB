package db

import (
	"log"
	"os"
	"strings"
	"time"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var LocalSession *gocql.Session
var AstraSession *gocql.Session

// GetSession returns the appropriate Cassandra session based on the mode and operation.
// opType: "read" or "write" (optional, for future logic)
func GetSession(opType ...string) *gocql.Session {
	mode := os.Getenv("CASSANDRA_MODE")
	switch mode {
	case "local":
		return LocalSession
	case "astra":
		return AstraSession
	default:
		return nil
	}
}

func InitCassandra() {
	_ = godotenv.Load() // Loads .env or .env.astra

	mode := os.Getenv("CASSANDRA_MODE") // local or astra

	// Local nodes
	if mode == "local" {
		hosts := os.Getenv("CASSANDRA_HOST")
		keyspace := os.Getenv("CASSANDRA_KEYSPACE")
		if hosts == "" || keyspace == "" {
			log.Fatal("Missing environment variables for local Cassandra")
		}
		hostList := strings.Split(hosts, ",")
		cluster := gocql.NewCluster(hostList...)
		cluster.Keyspace = keyspace
		// Use Consistency QUORUM for strong consistency: with 3 nodes and RF=3, at least 2 nodes must acknowledge
		cluster.Consistency = gocql.Quorum
		cluster.Timeout = 10 * time.Second

		maxRetries := 0 // 0 = infinito
		retryDelay := 10 * time.Second

		for attempts := 0; maxRetries == 0 || attempts < maxRetries; attempts++ {
			session, err := cluster.CreateSession()
			if err == nil {
				LocalSession = session
				break
			}
			log.Printf("Error connecting to Cassandra: %v. Retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
		}
		if LocalSession == nil {
			log.Fatal("Could not connect to any Cassandra node after retries")
		}
	}

	// Astra
	if mode == "astra" {
		dbID := os.Getenv("ASTRA_DB_ID")
		token := os.Getenv("APPLICATION_TOKEN")
		keyspace := os.Getenv("KEYSPACE_NAME")
		if dbID == "" || token == "" || keyspace == "" {
			log.Fatal("Missing environment variables for Astra DB")
		}
		cluster, err := gocqlastra.NewClusterFromURL(
			"https://api.astra.datastax.com",
			dbID,
			token,
			10*time.Second,
		)
		if err != nil {
			log.Fatalf("unable to load cluster from astra: %v", err)
		}
		cluster.Keyspace = keyspace
		cluster.Timeout = 30 * time.Second
		session, err := gocql.NewSession(*cluster)
		if err != nil {
			log.Fatalf("unable to connect session to Astra: %v", err)
		}
		AstraSession = session
	}
}
