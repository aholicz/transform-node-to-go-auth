package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"aholicz.lab.go.auth.jwt/auth"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/tOnkowzl/libs/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var router = echo.New()

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("fatal error config file: %s", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	runtime.GOMAXPROCS(viper.GetInt("app.maxprocs"))
}

func init() {
	router.HideBanner = true
	router.HidePort = true
}

func main() {
	dbConn := newDatabaseConnection()
	defer close(dbConn)

	authHander := auth.NewHandler(auth.NewService(auth.NewAuthDataSource(dbConn)))

	router.POST("/create-auth-profile", authHander.CreateAuthProfile)
	router.POST("/authenticate", authHander.Authenticate)

	go startServer()
	shutdown()
}

func startServer() {
	logx.Infof("start at port %s", viper.GetString("app.port"))
	logx.Info(router.Start(":" + viper.GetString("app.port")))
}

func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	if err := router.Shutdown(context.Background()); err != nil {
		logx.Panic("shutdown server:", err)
	}
}

func newDatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?checkConnLiveness=false&loc=Local&parseTime=true&readTimeout=%s&timeout=%s&writeTimeout=%s&maxAllowedPacket=0",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host")+":"+viper.GetString("mysql.port"),
		viper.GetString("mysql.database"),
		viper.GetString("mysql.timeout"),
		viper.GetString("mysql.timeout"),
		viper.GetString("mysql.timeout"),
	)

	logx.Errorf("[CONFIG] [MYSQL] %s", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logx.DefaultGormLogger})
	if err != nil {
		logx.Panicf("cannot open mysql connection:%s", err)
	}

	mysql, err := db.DB()
	if err != nil {
		logx.Panic(err)
	}

	mysql.SetMaxIdleConns(viper.GetInt("mysql.maxidle"))
	mysql.SetMaxOpenConns(viper.GetInt("mysql.maxidle"))
	mysql.SetConnMaxLifetime(viper.GetDuration("mysql.maxlifetime"))

	return db
}

func close(db *gorm.DB) {
	mysql, err := db.DB()
	if err != nil {
		logx.Errorf("can't get mysql from gorm cause %s", err)
		return
	}

	mysql.Close()
}
