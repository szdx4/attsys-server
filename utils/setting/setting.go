package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	// RunMode 表示应用当前运行模式
	RunMode string

	// HTTPPort 表示 HTTP 监听端口
	HTTPPort int

	// ReadTimeout 表示 HTTP Server 读取超时时间
	ReadTimeout time.Duration

	// WriteTimeout 表示 HTTP Server 写入超时时间
	WriteTimeout time.Duration
)

var cfg *ini.File

func init() {
	var err error

	cfg, err = ini.Load("config/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse config/app.ini: %v\n", err)
	}

	LoadApp()
	LoadServer()
}

// LoadApp 加载 APP 段配置
func LoadApp() {
	sec, err := cfg.GetSection("APP")
	if err != nil {
		log.Fatalf("Fail to get section APP in config/app.ini: %v\n", err)
	}

	RunMode = sec.Key("RUN_MODE").MustString("debug")
}

// LoadServer 加载 SERVER 段配置
func LoadServer() {
	sec, err := cfg.GetSection("SERVER")
	if err != nil {
		log.Fatalf("Fail to get section SERVER in config/app.ini: %v\n", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

// GetDatabase 获取 DATABASE 段配置
func GetDatabase() map[string]string {
	sec, err := cfg.GetSection("DATABASE")
	if err != nil {
		log.Fatalf("Fail to get section DATABASE in config/app.ini: %v\n", err)
	}

	return map[string]string{
		"type":     sec.Key("TYPE").String(),
		"host":     sec.Key("HOST").String(),
		"user":     sec.Key("USER").String(),
		"password": sec.Key("PASSWORD").String(),
		"name":     sec.Key("NAME").String(),
	}
}
