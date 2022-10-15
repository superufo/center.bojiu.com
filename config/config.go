package config

import (
	"math/rand"
	"time"

	"common.bojiu.com/discover/kit/sd/etcdv3"
	"google.golang.org/grpc"
)

var GlobalCfg = Config{}

type ServiceConfig struct {
	Storage struct {
		Protocol   string   `yaml:"protocol"`                             //协议
		Ip         string   `yaml:"ip"`                                   //ip
		Port       int      `yaml:"port"`                                 //端口
		EtcdServer []string `yaml:"etcdServer" mapstructure:"etcdServer"` //etcd配置
		RegKey     string   `yaml:"regKey"`                               //注册
	} `yaml:"storage" mapstructure:"storage"`
	Stream struct {
		Protocol string `yaml:"protocol"` //协议
		Ip       string `yaml:"ip"`       //ip
		Port     int    `yaml:"port"`     //端口
	} `yaml:"stream" mapstructure:"stream"`
}

type LogConfig struct {
	Path            string `yaml:"path"`            //日志输出路径
	Prefix          string `yaml:"prefix"`          //日志文件前缀
	Level           string `yaml:"level"`           //日志级别：debug/info/error/warn
	Development     bool   `yaml:"development"`     //是否为开发者模式
	DebugFileSuffix string `yaml:"debugFileSuffix"` //debug日志文件后缀
	WarnFileSuffix  string `yaml:"warnFileSuffix"`  //warn日志文件后缀
	ErrorFileSuffix string `yaml:"errorFileSuffix"` //error日志文件后缀
	InfoFileSuffix  string `yaml:"infoFileSuffix"`  //info日志文件后缀
	MaxAge          int    `yaml:"maxAge"`          //保存的最大天数
	MaxBackups      int    `yaml:"maxBackups"`      //最多存在多少个切片文件
	MaxSize         int    `yaml:"maxSize"`         //日日志文件大小（M）
}

type MysqlConfig struct {
	Addr     string `yaml:"addr"`     //地址
	Port     int    `yaml:"port"`     //端口
	Username string `yaml:"username"` //用户名
	Password string `yaml:"password"` //密码
	Database string `yaml:"database"` //库名
}

type RedisConfig struct {
	Addr         string `yaml:"addr"`         //地址:端口
	Password     string `yaml:"password"`     //密码
	DB           int    `yaml:"db"`           //库
	PoolSize     int    `yaml:"poolSize"`     //连接池大小
	MinIdleConns int    `yaml:"minIdleConns"` //最小空闲连接
}

type Config struct {
	Title       string        `yaml:"title"`                                    //服务器标题
	Active      string        `yaml:"active"`                                   //开发dev  测试test 上线 pro
	Service     ServiceConfig `yaml:"ser" mapstructure:"ser"`                   //服务器配置
	Log         LogConfig     `yaml:"log" mapstructure:"log"`                   //日志配置
	MysqlMaster MysqlConfig   `yaml:"mysql-master" mapstructure:"mysql-master"` //mysql主库
	MysqlLog    MysqlConfig   `yaml:"mysql-log" mapstructure:"mysql-log"`       //mysql日志库
	MysqlSlave1 MysqlConfig   `yaml:"mysql-slave1" mapstructure:"mysql-slave1"` //mysql从库
	Redis       RedisConfig   `yaml:"redis" mapstructure:"redis"`               //redis配置
}

type serverCfg struct {
	protocol   string
	ip         string
	port       int
	etcdServer []string
	regKey     string

	option etcdv3.ClientOptions
}

func NewServerCfg() *serverCfg {
	cfg := GlobalCfg.Service.Storage
	return &serverCfg{
		protocol:   cfg.Protocol,
		ip:         cfg.Ip,
		port:       cfg.Port,
		etcdServer: cfg.EtcdServer,
		regKey:     cfg.RegKey,
		option: etcdv3.ClientOptions{
			// Path to trusted ca file
			CACert: "",
			// Path to certificate
			Cert: "",
			// Path to private key
			Key: "",
			// Username if required
			Username: "",
			// Password if required
			Password: "",
			// If DialTimeout is 0, it defaults to 3s
			DialTimeout: time.Second * 3,
			// If DialKeepAlive is 0, it defaults to 3s
			DialKeepAlive: time.Second * 3,
			// If passing `grpc.WithBlock`, dial connection will block until success.
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
		},
	}
}

func (f *serverCfg) GetProtocol() string {
	return f.protocol
}

func (f *serverCfg) GetIp() string {
	return f.ip
}

func (f *serverCfg) GetPort() int {
	return f.port
}

func (f *serverCfg) GetEtcdServer() []string {
	return f.etcdServer
}

func (f *serverCfg) GetRandEtcdServer() string {
	n := rand.Intn(len(f.etcdServer))
	return f.etcdServer[n]
}

func (f *serverCfg) GetRegKey() string {
	return f.regKey
}

func (f *serverCfg) GetOption() etcdv3.ClientOptions {
	return f.option
}
