package db

import (
	"database/sql"
	"errors"
	"fmt"
	"service/pkg/log"
	"time"

	//github.com/jackc/pgx/v4/stdlib
	//_ "github.com/jackc/pgx/v4/stdlib"

	//github.com/lib/pq
	_ "github.com/lib/pq"
)

//Config database struct
type Config struct {
	Driver          string
	Host            string
	Port            string
	Dbname          string
	SslMode         string
	User            string
	Pass            string
	ConnMaxLifetime int
	MaxOpenConns    int
	MaxIdleConns    int
	ApplicationName string
}

var cfg *Config

var db *sql.DB
var err error

//ErrNoRows database/sql
var ErrNoRows = sql.ErrNoRows

//NullString from database/sql
type NullString = sql.NullString

//SetConfigDB ...
func SetConfigDB(conf *Config) {
	cfg = conf
}

//Init func
func Init() {

	if cfg == nil {
		log.Info("config is nil Error config Database not set", cfg)
		panic(errors.New("config is nil"))
	}

	dbConnString := ""
	if cfg.Driver == "mysql" {
		dbConnString = fmt.Sprintf("%s:%s@/%s", cfg.User, cfg.Pass, cfg.Dbname)

	} else {
		dbConnString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s application_name=%s",
			cfg.User,
			cfg.Pass,
			cfg.Host,
			cfg.Port,
			cfg.Dbname,
			cfg.SslMode,
			cfg.ApplicationName,
		)
	}

	db, err = sql.Open("postgres", dbConnString)

	if err != nil {
		log.Error("Failed to connect to database >> ", dbConnString, err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	 	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	   	db.SetMaxIdleConns(cfg.MaxIdleConns)

	   	// SetMaxOpenConns sets the maximum number of open connections to the database.
	   	db.SetMaxOpenConns(cfg.MaxOpenConns)

	   	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	   	db.SetConnMaxLifetime(time.Millisecond * time.Duration(cfg.ConnMaxLifetime)) 

	

}

//GetDB - get DB
func GetDB() *sql.DB {
	/* err = db.Ping()
	if err != nil {
		panic(err)
	} */
	return db
}

/*
+
+
+
+		SafecityDB
+
+
+
*/

var scfg *Config

var sdb *sql.DB

//SetConfigDBS ...
func SetConfigDBS(conf *Config) {
	scfg = conf
}

//InitS ....
func InitS() {

	if scfg == nil {
		log.Info("config is nil Error config Database not set", cfg)
		panic(errors.New("config is nil"))
	}

	dbConnString := ""

	dbConnString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s application_name=%s",
		scfg.User,
		scfg.Pass,
		scfg.Host,
		scfg.Port,
		scfg.Dbname,
		scfg.SslMode,
		scfg.ApplicationName,
	)

	sdb, err = sql.Open("postgres", dbConnString)

	if err != nil {
		log.Error("Failed to connect to database >> ", dbConnString, err)
		panic(err)
	}

	err = sdb.Ping()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sdb.SetMaxIdleConns(scfg.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sdb.SetMaxOpenConns(scfg.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sdb.SetConnMaxLifetime(time.Millisecond * time.Duration(scfg.ConnMaxLifetime))

}

//GetDBS - get DBS
func GetDBS() *sql.DB {
	/* err = sdb.Ping()
	if err != nil {
		log.Error("get DBS", scfg, err)
		panic(err)
	} */
	return sdb
}
