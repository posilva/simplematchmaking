package config

import (
	"github.com/spf13/viper"
)

const (
	httpAddr  = "ADDR"
	redisAddr = "REDIS_ADDR"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("SLMM")
	viper.SetConfigName("simplematchmaking")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		viper.AutomaticEnv()
	}

	// set defaults
	viper.SetDefault(httpAddr, ":8808")
	viper.SetDefault(redisAddr, "localhost:6379")

}

// GetAddr returns the http server addresss
func GetAddr() string {
	return viper.GetString(httpAddr)
}

func GetRedisAddr() string {
	return viper.GetString(redisAddr)
}

func IsLocal() bool {
	return viper.GetBool("local")
}

func SetRedisAddr(v string) {
	viper.Set(redisAddr, v)
}
func SetAddr(v string) {
	viper.Set(httpAddr, v)
}
func SetLocal(v bool) {
	viper.Set("local", v)
}
