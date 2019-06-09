package config

import (
	"log"
	"time"

	"github.com/Unknwon/goconfig"
)

// AppConfig 存储应用配置
type AppConfig struct {
	RunMode            string
	EncryptKey         string
	MinPwdLength       int
	ItemsPerPage       int
	TokenValid         int
	MinUserLength      int
	MinOvertimeMinutes int
}

// ServerConfig 存储 HTTP 服务器配置
type ServerConfig struct {
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig 存储数据库连接信息
type DatabaseConfig struct {
	Type     string
	Host     string
	User     string
	Password string
	Name     string
}

// ConnectionString 返回数据库连接信息字符串
func (c *DatabaseConfig) ConnectionString() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ")/" + c.Name + "?charset=utf8&parseTime=True&loc=Asia%2fShanghai"
}

// QcloudConfig 存储腾讯云 API 信息
type QcloudConfig struct {
	SecretID  string
	SecretKey string
}

var cfg *goconfig.ConfigFile

// 读取的配置信息
var (
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Qcloud   QcloudConfig
)

func init() {
	var err error
	cfg, err = goconfig.LoadConfigFile("./config/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse config/app.ini: %v", err)
	}

	loadApp()
	loadServer()
	loadDatabase()
	loadQcloud()
}

func loadApp() {
	App.RunMode = cfg.MustValue("APP", "RUN_MODE", "release")
	App.EncryptKey = cfg.MustValue("APP", "ENCRYPT_KEY", "attsys-server")
	App.MinPwdLength = cfg.MustInt("APP", "MIN_PWD_LENGTH", 7)
	App.ItemsPerPage = cfg.MustInt("APP", "ITEMS_PER_PAGE", 20)
	App.TokenValid = cfg.MustInt("APP", "TOKEN_VALID", 2)
	App.MinUserLength = cfg.MustInt("APP", "MIN_USER_LENGTH", 2)
	App.MinOvertimeMinutes = cfg.MustInt("APP", "MIN_OVERTIME_MINUTES", 60)
}

func loadServer() {
	Server.HTTPPort = cfg.MustInt("SERVER", "HTTP_PORT", 8000)
	Server.ReadTimeout = time.Duration(cfg.MustInt("SERVER", "READ_TIMEOUT", 60)) * time.Second
	Server.WriteTimeout = time.Duration(cfg.MustInt("SERVER", "WRITE_TIMEOUT", 60)) * time.Second
}

func loadDatabase() {
	Database.Type = cfg.MustValue("DATABASE", "TYPE", "mysql")
	Database.Host = cfg.MustValue("DATABASE", "HOST", "127.0.0.1:3306")
	Database.User = cfg.MustValue("DATABASE", "USER", "root")
	Database.Password = cfg.MustValue("DATABASE", "PASSWORD", "root")
	Database.Name = cfg.MustValue("DATABASE", "NAME", "attsys")
}

func loadQcloud() {
	Qcloud.SecretID = cfg.MustValue("QCLOUD", "SECRET_ID", "")
	Qcloud.SecretKey = cfg.MustValue("QCLOUD", "SECRET_KEY", "")
}

// StrToTime 将接收到的时间字符串转换成 time.Time
func StrToTime(str string) (startAt time.Time, err error) {
	var timeLayOut = "2006-01-02 15:04:05"
	startAt, err = time.Parse(timeLayOut, str)
	return
}
