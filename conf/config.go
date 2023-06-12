package conf

import (
	"github.com/spf13/viper"
)

func Init(confName string) error {
	viper.AddConfigPath("./conf")
	viper.SetConfigName(confName)

	return viper.ReadInConfig()
}
