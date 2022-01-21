package config

import (
	"github.com/spf13/viper"
	"strings"
)

type ViperLoader struct{}

func (v *ViperLoader) Load(params ...interface{}) (*Config, error) {
	//viper.SetConfigName("config")
	viper.SetEnvPrefix("SLOW")
	//viper.SetConfigType("json")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//if err := viper.ReadInConfig(); err != nil {
	//	return nil, err
	//}
	c := &Config{
		AppAddress: viper.GetString("app.address"),
		DBType:     convertStrToDB(viper.GetString("db.type")),
		AppPort:    viper.GetInt("app.port"),
		LogFile:    viper.GetString("db.logfile"),
	}

	return c, nil
}

// convertStrToDB converts string to DBType. It returns postgres if there isn't any match
func convertStrToDB(s string) DBType {
	s = strings.ToLower(s)
	if s == "mysql" {
		return DBMySql
	}
	return DBPostgres
}
