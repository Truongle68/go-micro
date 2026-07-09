package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	_defaultMaxPoolSize = 10
	_defaultConnAttempt = 3
	_defaultConnTimeout = time.Second
)

type Mongo struct {
	maxPoolSize int
	connAttempt int
	connTimeout time.Duration

	Client   *mongo.Client
	Database *mongo.Database
}

func New(url, dbName string, opts ...Option) (*Mongo, error) {
	mg := &Mongo{
		maxPoolSize: _defaultMaxPoolSize,
		connAttempt: _defaultConnAttempt,
		connTimeout: _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(mg)
	}

	clientOpts := options.Client().ApplyURI(url).SetMaxPoolSize(uint64(mg.maxPoolSize))

	var err error

	for mg.connAttempt > 0 {
		mg.Client, err = mongo.Connect(clientOpts)

		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), mg.connTimeout)
			err = mg.Client.Ping(ctx, readpref.Primary())
			cancel()
		}

		if err == nil {
			break
		}

		log.Printf("MongoDB is trying to connect, attempts left: %d", mg.connAttempt)

		time.Sleep(mg.connTimeout)

		mg.connAttempt--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	mg.Database = mg.Client.Database(dbName)

	return mg, nil
}

func (m *Mongo) Close() {
	if m.Client != nil {
		_ = m.Client.Disconnect(context.Background())
	}
}
