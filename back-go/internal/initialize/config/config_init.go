package config

import (
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../internal/initialize/config")   //开发环境
	viper.AddConfigPath("/app/internal/initialize/config") //运行环境
	viper.AddConfigPath("./internal/initialize/config")    //容器环境
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		return err
	}
	return nil
}

var Conf = Config{}

type Config struct {
	APP struct {
		Name     string
		IP       string
		Port     int
		Mode     string
		StaticFS bool
	}
	MySQL struct {
		DriverName string `yaml:"driverName"`
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		Database   string `yaml:"database"`
		UserName   string `yaml:"username"`
		Password   string `yaml:"password"`
		Charset    string `yaml:"charset"`
	}
	Redis struct {
		Addr         string `yaml:"addr"`
		Port         string `yaml:"port"`
		Password     string `yaml:"password"`
		DB           int    `yaml:"db"`
		PoolSize     int    `yaml:"poolSize"`
		MinIdleConns int    `yaml:"minIdleConns"`
	}
}
