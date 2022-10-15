package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Vp = viper.New()
)

func InitVp() {
	Vp.SetConfigName("config/config")
	Vp.AddConfigPath("../")
	Vp.AddConfigPath(".") // 添加搜索路径
	err := Vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s", err))
	}
	err = Vp.Unmarshal(&GlobalCfg)
	if err != nil {
		panic(fmt.Errorf("配置文件解析失败: %s", err))
	}
	Vp.WatchConfig()
	Vp.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件变更")
		err := Vp.Unmarshal(&GlobalCfg)
		if err != nil {
			fmt.Println("配置文件更新失败")
		} else {
			fmt.Printf("配置文件更新成功: %+v \n", GlobalCfg)
		}
	})
}
