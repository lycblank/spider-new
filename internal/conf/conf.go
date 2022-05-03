package conf

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	FlyBook FlyBookConfig `yaml:"flybook"`
	Chanify ChanifyConfig `yaml:"chanify"`
	PushPlus PushPlusConfig `yaml:"pushplus"`
	DingDing DingDingConfig `yaml:"dingding"`
}

type DingDingConfig struct {
	Webhook string `yaml:"webhook"`
	Secret string `yaml:"secret"`
}

type ChanifyConfig struct {
	Webhook string `yaml:"webhook"`
}

type FlyBookConfig struct {
	Webhook string `yaml:"webhook"`
}

type PushPlusConfig struct {
	Webhook string  `yaml:"webhook"`
	Group string    `yaml:"group"`
	Token string    `yaml:"token"`
}

var config *Config
var configOnce sync.Once
func GetConfig() *Config {
	configOnce.Do(func(){
		if err := godotenv.Load(); err != nil {
			fmt.Println("not found .env file")
		}
		viper.AutomaticEnv()
		confPath := viper.GetString("CONF_PATH")
		if confPath == "" {
			confPath = "configs"
		}
		viper.SetConfigName("service")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(confPath)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		cfg := &Config{}
		if err := viper.Unmarshal(cfg, func(v *mapstructure.DecoderConfig){
			v.TagName = "yaml"
		}); err != nil {
			panic(err)
		}
		config = cfg
	})
	return config
}



