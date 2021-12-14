package orm

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
)

type Config struct {
	WriteTimeout    string `yaml:"writeTimeout"`
	Loc             string `yaml:"loc"`
	Port            int    `yaml:"port"`
	ReadTimeout     string `yaml:"readTimeout"`
	Charset         string `yaml:"charset"`
	ParseTime       bool   `yaml:"parseTime"`
	Protocol        string `yaml:"protocol"`
	Dsn             string `yaml:"dsn"`
	Database        string `yaml:"database"`
	Collation       string `yaml:"collation"`
	Timeout         string `yaml:"timeout"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	ConnMaxIdle     int    `yaml:"connMaxIdle"`
	ConnMaxOpen     int    `yaml:"connMaxOpen"`
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
}

func GetBaseConfig(c framework.Container) *Config {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &Config{}
	err := configService.Load("database", config)
	if err != nil {
		logService.Error(context.Background(), "parse database config error", nil)
		return nil
	}
	return config
}

// FormatDsn 生成dsn
func (conf *Config) FormatDsn() (string, error) {
	port := strconv.Itoa(conf.Port)
	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		fmt.Println("database timeout parse error :", err.Error())
		return "", err
	}
	readTimeout, err := time.ParseDuration(conf.ReadTimeout)
	if err != nil {
		fmt.Println("database readTimeout parse error :", err.Error())
		return "", err
	}
	writeTimeout, err := time.ParseDuration(conf.WriteTimeout)
	if err != nil {
		fmt.Println("database writeTimeout parse error :", err.Error())
		return "", err
	}
	location, err := time.LoadLocation(conf.Loc)
	if err != nil {
		fmt.Println("database localtion parse error :", err.Error())
		return "", err
	}
	driverConf := &mysql.Config{
		User:         conf.Username,
		Passwd:       conf.Password,
		Net:          conf.Protocol,
		Addr:         net.JoinHostPort(conf.Host, port),
		DBName:       conf.Database,
		Collation:    conf.Collation,
		Loc:          location,
		Timeout:      timeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ParseTime:    conf.ParseTime,
	}
	return driverConf.FormatDSN(), nil
}
