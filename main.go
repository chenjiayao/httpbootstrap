package main

import (
	"flag"

	"github.com/spf13/viper"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "./config/development.ini", "指定配置文件")
	flag.Parse()
	loadConfig(configFile)

	run()
}

func loadConfig(filename string) {
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
