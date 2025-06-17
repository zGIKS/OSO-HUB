package db

import (
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var Session *gocql.Session

func InitCassandra() {
	_ = godotenv.Load() // Carga .env si existe
	host := os.Getenv("CASSANDRA_HOST")
	port := os.Getenv("CASSANDRA_PORT")
	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	var cluster *gocql.ClusterConfig
	if port != "" {
		cluster = gocql.NewCluster(host + ":" + port)
	} else {
		cluster = gocql.NewCluster(host)
	}
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Error connecting to Cassandra: %v", err)
	}
	Session = session
}
