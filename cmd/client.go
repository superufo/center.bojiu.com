package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"center.bojiu.com/internal/net/storage/proto"
// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"

// 	"center.bojiu.com/pkg/log"

// 	"center.bojiu.com/config"
// )

// func testClient() {
// 	// 初始化配置文件
// 	config.InitVp()

// 	// 初始化日志文件
// 	log.ZapLog = log.InitLogger()

// 	scfg := config.NewServerCfg()
// 	log.ZapLog.Info(fmt.Sprintf("%s:%d", scfg.GetIp(), scfg.GetPort()))

// 	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", scfg.GetIp(), scfg.GetPort()), grpc.WithInsecure())

// 	if err != nil {
// 		log.ZapLog.With(zap.Error(err)).Error("grpc dial error")
// 	}
// 	defer conn.Close()

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client := proto.NewStorageClient(conn)
// 	request := proto.GameSettlementTos{
// 		GameType: 11,
// 		Score:    120,
// 	}
// 	reply, err := client.SendSettlement(ctx, &request)

// 	if err != nil {
// 		log.ZapLog.With(zap.Error(err)).Error("grpc dial result")
// 	}
// 	log.ZapLog.With(zap.Any("reply", reply)).Info("grpc result")

// }
