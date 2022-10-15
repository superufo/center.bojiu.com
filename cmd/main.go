package main

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"center.bojiu.com/config"
	storage "center.bojiu.com/internal/net/storage"
	proto "center.bojiu.com/internal/net/storage/proto"
	"center.bojiu.com/pkg/log"
	"center.bojiu.com/pkg/mysql"
	"center.bojiu.com/pkg/redislib"
	etcdv3 "common.bojiu.com/discover/kit/sd/etcdv3"
)

func main() {
	// 初始化配置文件
	config.InitVp()
	// 初始化日志文件
	log.ZapLog = log.InitLogger()
	// 初始化redis
	redis := redislib.GetClient()
	defer redis.Close()
	// 初始化数据库 获取 mysql.M()  mysql.S()
	MasterDB := mysql.MasterInit()
	defer MasterDB.Close()
	LogDB := mysql.LogDBInit()
	defer LogDB.Close()
	Slave1DB := mysql.Slave1Init()
	defer Slave1DB.Close()

	scfg := config.NewServerCfg()
	grpcServer := grpc.NewServer()
	storageServer := storage.NewStorageServerImpl()
	proto.RegisterStorageServer(grpcServer, storageServer)

	//开启定时任务
	stopChan := make(chan struct{})
	defer close(stopChan)
	storageServer.RunTask(stopChan)
	//监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", scfg.GetPort()))
	if err != nil {
		log.ZapLog.With(zap.Any("err", err), zap.Stack("trace")).Info("grpc")
	}
	/*******服务注册 start*******/
	instance := fmt.Sprintf("%s://%s:%d", scfg.GetProtocol(), scfg.GetIp(), scfg.GetPort())
	client, err := etcdv3.NewClient(context.Background(), scfg.GetEtcdServer(), scfg.GetOption())
	if err != nil {
		log.ZapLog.With(zap.Error(err), zap.Stack("trace")).Info("error")
		return
	}
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{
		Key:   scfg.GetRegKey(),
		Value: instance,
	}, log.ZapLog)

	registrar.Register()
	defer registrar.Deregister()
	v, _ := client.GetEntries(scfg.GetRegKey())
	log.ZapLog.With(zap.Any("regKey", v), zap.Stack("trace")).Info("main")
	/*******服务注册 end *******/

	log.ZapLog.Info("中心服务器启动........")
	grpcServer.Serve(lis)
}
