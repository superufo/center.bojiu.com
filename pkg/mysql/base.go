package mysql

import "github.com/go-xorm/xorm"

func M() *xorm.Engine {
	return MasterDB
}

func L() *xorm.Engine {
	return LogDB
}

func S1() *xorm.Engine {
	return Slave1DB
}

//oxrm DB事务封装
type DBTransFunc func(session *xorm.Session) error

func RunTrans(engine *xorm.Engine, f DBTransFunc) error {
	var err error
	s := engine.NewSession()
	defer s.Close()
	if err = s.Begin(); err != nil {
		return err
	}
	err = f(s) //将session传入回调函数执行sql操作，此回调函数将返回err，如果sql操作有错则err不为空，将会在commit之前实现回滚
	if err != nil {
		s.Rollback()
		return err
	} else if err = s.Commit(); err != nil {
		return err
	}
	return nil
}
