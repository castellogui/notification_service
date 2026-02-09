package infra

import (
	"time"

	"github.com/gocql/gocql"
)

type ScyllaConfig struct {
	Hosts       []string
	Port        int
	Keyspace    string
	Username    string
	Password    string
	Consistency gocql.Consistency
}

func NewScyllaSession(cfg ScyllaConfig) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	if cfg.Port != 0 {
		cluster.Port = cfg.Port
	}
	cluster.Keyspace = cfg.Keyspace
	if cfg.Consistency != 0 {
		cluster.Consistency = cfg.Consistency
	} else {
		cluster.Consistency = gocql.LocalQuorum
	}

	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second
	cluster.NumConns = 2 // per host; tune later

	if cfg.Username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	// TLS if needed:
	// cluster.SslOpts = &gocql.SslOptions{Config: tlsConfig}

	return cluster.CreateSession()
}