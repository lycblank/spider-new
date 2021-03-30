package conf

import (
	"sync"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	FlyBook FlyBookConfig `yaml:"flybook"`
}

type FlyBookConfig struct {
	Webhook string `yaml:"webhook"`
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



