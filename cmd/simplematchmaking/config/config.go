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
	viper.SetEnvPrefix("SMM")
	viper.SetConfigName("simplematchmaking")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		viper.AutomaticEnv()
	}

	// set defaults
	viper.SetDefault(httpAddr, ":8808")
	viper.SetDefault(redisAddr, "localhost:6379")

	viper.BindEnv(httpAddr)
	viper.BindEnv(redisAddr)

}

// GetAddr returns the http server addresss
func GetAddr() string {
	return viper.GetString(httpAddr)
}

// GetRedisAddr returns the redis server address
func GetRedisAddr() string {
	return viper.GetString(redisAddr)
}

// IsLocal returns true if the server is running in local mode
func IsLocal() bool {
	return viper.GetBool("local")
}

// SetRedisAddr sets the redis server address
func SetRedisAddr(v string) {
	viper.Set(redisAddr, v)
}

// SetAddr sets the http server address
func SetAddr(v string) {
	viper.Set(httpAddr, v)
}

// SetLocal sets the local mode
func SetLocal(v bool) {
	viper.Set("local", v)
}
