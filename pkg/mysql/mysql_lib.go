package mysql

import (
	"fmt"

	"time"

	"center.bojiu.com/config"
	"center.bojiu.com/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

var (
	MasterDB *xorm.Engine
	LogDB    *xorm.Engine
	Slave1DB *xorm.Engine
)

func MasterInit() *xorm.Engine {
	cfg := config.GlobalCfg
	master := cfg.MysqlMaster
	open := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		master.Username,
		master.Password,
		master.Addr,
		master.Port,
		master.Database)
	var err error
	MasterDB, err = xorm.NewEngine("mysql", open)
	if err != nil {
		log.ZapLog.Error(fmt.Sprintf("Open mysql-master failed,err:%s\n", err.Error()))
		panic(err)
	}
	MasterDB.SetTableMapper(core.SameMapper{})
	MasterDB.SetColumnMapper(core.GonicMapper{})
	MasterDB.SetConnMaxLifetime(100 * time.Second)
	MasterDB.SetMaxOpenConns(100)
	MasterDB.SetMaxIdleConns(16)
	err = MasterDB.Ping()
	if err != nil {
		log.ZapLog.Error(fmt.Sprintf("Failed to connect to mysql-master, err:%s" + err.Error()))
		panic(err.Error())
	}
	// 显示打印语句
	if cfg.Active == "dev" || cfg.Active == "test" {
		MasterDB.ShowSQL(true)
	}
	//MasterDB.Sync2(new(bj_server.GamesInfo))
	log.ZapLog.Info("mysql-master connect success\r\n")
	return MasterDB
}

func LogDBInit() *xorm.Engine {
	cfg := config.GlobalCfg
	l := cfg.MysqlLog
	open := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		l.Username,
		l.Password,
		l.Addr,
		l.Port,
		l.Database)

	var err error
	LogDB, err = xorm.NewEngine("mysql", open)
	if err != nil {
		log.ZapLog.Error(fmt.Sprintf("Open mysql-master failed,err:%s\n", err.Error()))
		panic(err)
	}
	LogDB.SetTableMapper(core.SameMapper{})
	LogDB.SetColumnMapper(core.GonicMapper{})
	LogDB.SetConnMaxLifetime(100 * time.Second)
	LogDB.SetMaxOpenConns(100)
	LogDB.SetMaxIdleConns(16)
	err = LogDB.Ping()
	if err != nil {
		log.ZapLog.Error(fmt.Sprintf("Failed to connect to mysql-master, err:%s" + err.Error()))
		panic(err.Error())
	}
	// 显示打印语句
	if cfg.Active == "dev" || cfg.Active == "test" {
		LogDB.ShowSQL(true)
	}
	log.ZapLog.Info("mysql-master connect success\r\n")
	return LogDB
}

func Slave1Init() *xorm.Engine {
	cfg := config.GlobalCfg
	Slave1 := cfg.MysqlSlave1

	open := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		Slave1.Username,
		Slave1.Password,
		Slave1.Addr,
		Slave1.Port,
		Slave1.Database)
	var err error
	Slave1DB, err = xorm.NewEngine("mysql", open)
	if err != nil {
		log.ZapLog.Error(fmt.Sprintf("Open mysql-slave1 failed,err:%s\n", err.Error()))
		panic(err)
	}
	Slave1DB.SetConnMaxLifetime(100 * time.Second)
	Slave1DB.SetMaxOpenConns(100)
	Slave1DB.SetMaxIdleConns(16)
	// 显示打印语句
	if cfg.Active == "dev" || cfg.Active == "test" {
		Slave1DB.ShowSQL(true)
	}
	err = Slave1DB.Ping()
	if err != nil {
		log.ZapLog.Error(fmt.Sprintf("Failed to connect to mysql-slave1, err:%s" + err.Error()))
		panic(err.Error())
	}
	log.ZapLog.Info("mysql-slave1 connect success\r\n")
	return Slave1DB
}
