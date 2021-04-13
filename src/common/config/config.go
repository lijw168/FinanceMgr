package config

import (
	"database/sql"
	"fmt"
	"time"

	"common/log"

	_ "github.com/go-sql-driver/mysql"
)

type ConfigCheck interface {
	CheckValid() error
}

type MysqlConf struct {
	Ip            string
	Port          int
	User          string
	Passwd        string
	DB            string
	Timeout       int
	MaxConnection int
	MaxLifetime   int
}

func (c *MysqlConf) CheckValid() error {
	if len(c.Ip) == 0 {
		return fmt.Errorf("invalid Ip")
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if len(c.User) == 0 || len(c.Passwd) == 0 || len(c.DB) == 0 {
		return fmt.Errorf("invalid MysqlConf")
	}
	if c.Timeout == 0 {
		c.Timeout = 3000
	}
	if c.MaxConnection == 0 {
		c.MaxConnection = 300
	}
	if c.MaxLifetime == 0 {
		c.MaxLifetime = 1000
	}
	return nil
}

type MysqlInstance struct {
	Conf   *MysqlConf
	Logger log.ILog
}

func (ins MysqlInstance) NewMysqlInstance() (*sql.DB, error) {
	strConn := "%s:%s@tcp(%s:%d)/%s?autocommit=true&parseTime=true&timeout=%dms&loc=Asia%%2FShanghai"
	url := fmt.Sprintf(strConn, ins.Conf.User, ins.Conf.Passwd,
		ins.Conf.Ip, ins.Conf.Port, ins.Conf.DB, ins.Conf.Timeout)
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", url)
	if err != nil {
		ins.Logger.Error("mysql open err: %s", err.Error())
		return nil, err
	}
	ins.Logger.Info("open mysql success\n")
	db.SetMaxOpenConns(ins.Conf.MaxConnection)
	db.SetMaxIdleConns(ins.Conf.MaxConnection)
	db.SetConnMaxLifetime(time.Second * time.Duration(ins.Conf.MaxLifetime))

	err = db.Ping()
	if err != nil {
		ins.Logger.Error("mysql ping err(%s)\n", err.Error())
		return nil, err
	}

	ins.Logger.Debug("[db] MySQLInit, configure: %+v", ins.Conf)
	return db, nil
}

type LogConf struct {
	Level       int    `json:"Level"`
	FileName    string `json:"FileName"`
	FileMaxSize int    `json:"FileMaxSize"`
	FileCount   int    `json:"FileCount"`
}

func (c *LogConf) CheckValid() error {
	if len(c.FileName) == 0 {
		return fmt.Errorf("AdminServerCfg need MysqlConf")
	}
	if c.FileMaxSize <= 1024 {
		c.FileMaxSize = 102400000
	}
	if c.FileCount < 1 {
		c.FileCount = 1
	}
	return nil
}

type LogFac struct {
	Logconf *LogConf
}

func (fac LogFac) NewLogger() (*log.Logger, error) {
	var h log.Handler
	var err error
	h, err = log.NewRotatingFileHandler(fac.Logconf.FileName, fac.Logconf.FileMaxSize, fac.Logconf.FileCount)
	if err != nil {
		fmt.Printf("new log handler err: %v\n", err.Error())
		return nil, err
	}

	logger := log.NewDefault(h)
	logger.SetLevel(fac.Logconf.Level)

	return logger, nil
}
