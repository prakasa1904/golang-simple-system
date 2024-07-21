package qr

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Repository struct {
	config *viper.Viper
	log    *logrus.Logger
	qr     chan string
}

func NewRepository(config *viper.Viper, log *logrus.Logger) *Repository {
	qr := make(chan string, 1)
	return &Repository{
		config: config,
		log:    log,
		qr:     qr,
	}
}
