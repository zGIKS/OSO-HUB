from cassandra.cluster import Cluster
from cassandra.auth import PlainTextAuthProvider
from dotenv import load_dotenv
import os

load_dotenv()

class CassandraConnector:
    def __init__(self):
        self.contact_points = [os.getenv("CASSANDRA_HOST", "127.0.0.1")]
        self.keyspace = os.getenv("CASSANDRA_KEYSPACE", "osohub")
        self.username = os.getenv("CASSANDRA_USER")
        self.password = os.getenv("CASSANDRA_PASS")
        self.session = None
        self.cluster = None

    def connect(self):
        if self.username and self.password:
            auth_provider = PlainTextAuthProvider(username=self.username, password=self.password)
            self.cluster = Cluster(self.contact_points, auth_provider=auth_provider)
        else:
            self.cluster = Cluster(self.contact_points)
        self.session = self.cluster.connect(self.keyspace)
        return self.session

    def shutdown(self):
        if self.session:
            self.session.shutdown()
        if self.cluster:
            self.cluster.shutdown()
