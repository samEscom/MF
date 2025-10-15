package db

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/samEscom/MF/users/src/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	maxIdle  = 1
	logQuery = true
	singular = true
)

type DBProvider int

var DBProviderNames = [...]string{
	"postgres",
	"sqlite3",
}

const (
	Postgresql DBProvider = iota
	SQLLite3
)

func (value DBProvider) String() string {
	return DBProviderNames[value]
}

type Connection struct {
	User        string
	Password    string
	Host        string
	Port        int
	DBName      string
	SSLMode     bool
	DBProvider  DBProvider
	PoolSize    int
	MaxIdleTime time.Duration
}

type AdditionalConfig struct {
	LogQuery      bool
	SingularTable bool
}

func (opt *Connection) SetDefaultValues() {
	const maxIdleSec = 60

	if opt.PoolSize == 0 {
		opt.PoolSize = 10
	}

	if opt.MaxIdleTime == 0 {
		opt.MaxIdleTime = maxIdleSec * time.Second
	}
}

func NewPostgresConnection() *gorm.DB {
	postgresCfg := config.Config().Db

	return CreateDataBaseConnection(Connection{
		User:        postgresCfg.User,
		Password:    postgresCfg.Password,
		Host:        postgresCfg.Host,
		Port:        postgresCfg.Port,
		DBName:      postgresCfg.DB,
		DBProvider:  Postgresql,
		PoolSize:    postgresCfg.PoolSize,
		MaxIdleTime: maxIdle,
	}, AdditionalConfig{
		LogQuery:      logQuery,
		SingularTable: singular,
	})
}

var (
	DBInstance  *gorm.DB
	onceElement sync.Once
)

func CreateDataBaseConnection(connectionParams Connection, additionalConfig AdditionalConfig) *gorm.DB {
	onceElement.Do(func() {
		DBInstance = getConnection(connectionParams, additionalConfig)
	})

	return DBInstance
}

func getConnection(connectionParams Connection, additionalConfig AdditionalConfig) *gorm.DB {
	connectionParams.SetDefaultValues()
	urlConnection, err := BuildURLConnection(connectionParams)

	if err != nil {
		log.Fatalf("error creating urlConnection for  %s database %s", connectionParams.DBProvider.String(), err)
		return nil
	}

	gormConfig := gorm.Config{}

	if additionalConfig.LogQuery {
		log.Print("enabling query log")

		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	if additionalConfig.SingularTable {
		log.Print("enabling singular tables")

		gormConfig.NamingStrategy = schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
			NameReplacer:  nil,
			NoLowerCase:   false,
		}
	}

	dialect, err := getDialect(connectionParams, urlConnection)
	if err != nil {
		log.Fatalf("error creating dialector for  %s database %s", connectionParams.DBProvider.String(), err)
		return nil
	}

	dbConnection, err := gorm.Open(dialect, &gormConfig)

	if err != nil {
		log.Fatalf("error on %s connection: %s",
			connectionParams.DBProvider.String(), err.Error())
		return nil
	}

	sqlDB, err := dbConnection.DB()
	if err != nil {
		log.Panicf("error on %s connection: %s",
			connectionParams.DBProvider.String(), err.Error())
	}

	sqlDB.SetConnMaxIdleTime(connectionParams.MaxIdleTime)
	sqlDB.SetMaxOpenConns(connectionParams.PoolSize)

	if err := sqlDB.Ping(); err != nil {
		log.Panicf("error on %s connection: %s",
			connectionParams.DBProvider.String(), err.Error())
	}

	log.Printf("%s successfully connected %s", connectionParams.DBProvider.String(), connectionParams.Host)

	return dbConnection
}

func getDialect(connectionParams Connection, urlConnection string) (gorm.Dialector, error) {
	switch connectionParams.DBProvider {
	case Postgresql:
		return postgres.Open(urlConnection), nil
	case SQLLite3:
		return sqlite.Open(urlConnection), nil
	default:
		return nil, errors.New("db provider unsupported")
	}
}

func BuildURLConnection(connectionParams Connection) (string, error) {
	SSLMode := validateSSLMode(connectionParams)

	switch connectionParams.DBProvider {
	case Postgresql:
		return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
			connectionParams.User,
			connectionParams.Password,
			connectionParams.Host,
			connectionParams.Port,
			connectionParams.DBName,
			SSLMode), nil
	case SQLLite3:
		return fmt.Sprintf("%s:memory:", connectionParams.DBName), nil
	default:
		return "", errors.New("db provider unsupported")
	}
}

func validateSSLMode(connectionParams Connection) string {
	if connectionParams.SSLMode {
		return "enable"
	} else {
		return "disable"
	}
}
