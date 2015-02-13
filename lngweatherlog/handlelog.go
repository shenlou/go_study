package main

import (
	"bufio"
	"github.com/Unknwon/goconfig"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigInfo struct {
	logPath string
	mailConfig
	timeConfig
}

type mailConfig struct {
	username, password, host, to string
}

type timeConfig struct {
	from, to string
}

//读取配置文件
func getLogFromIni() *ConfigInfo {
	configFile := "config.ini"
	key := "logPath"
	section := "mail"
	config, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Fatalf("载入配置文件%s错误,%s", configFile, err.Error())
	}
	logPath, err := config.GetValue(goconfig.DEFAULT_SECTION, key)
	if err != nil {
		log.Fatalf("读取节点%s错误,%s", key, err.Error())
	}
	username, err := config.GetValue(section, "username")
	if err != nil {
		log.Fatalf("读取节点%s错误%s", "username", err.Error())
	}
	password, err := config.GetValue(section, "password")
	if err != nil {
		log.Fatalf("读取节点%s错误%s", "password", err.Error())
	}
	host, err := config.GetValue(section, "host")
	if err != nil {
		log.Fatalf("读取节点%s错误%s", "host", err.Error())
	}
	to, err := config.GetValue(section, "to")
	if err != nil {
		log.Fatalf("读取节点%s错误%s", "to", err.Error())
	}
	timeFrom, err := config.GetValue("time", "from")
	if err != nil {
		log.Fatalf("读取节点%s错误%s,", "from", err.Error())
	}
	timeTo, err := config.GetValue("time", "to")
	if err != nil {
		log.Fatalf("读取节点%s错误%s,", "time.to", err.Error())
	}

	return &ConfigInfo{
		logPath:    logPath,
		mailConfig: mailConfig{username, password, host, to},
		timeConfig: timeConfig{timeFrom, timeTo},
	}
}

//处理lngweather日志
func handleLngweatherLog() {
	config := getLogFromIni()
	file, err := os.Open(config.logPath)
	if err != nil {
		log.Fatalf("读入文件%s出错,%s", config.logPath, err.Error())
	}
	defer file.Close()
	buff := bufio.NewReader(file)
	decode := mahonia.NewDecoder("gbk")
	nowDate := time.Now().Format("2006-01-02")
	successMsg := "LNG青岛接收站" + nowDate + "天气预报邮件发送成功！"
	a := 0
	for {
		line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		line = decode.ConvertString(line)
		if strings.Count(line, successMsg) == 1 {
			a += 1
		}
	}
	if a == 0 {
		errMsg := "LNG青岛接收站" + nowDate + "天气预报邮件发送失败！"
		log.Println(errMsg)
		sendLngErrLog(config.username, config.password, config.host, config.mailConfig.to, errMsg)
	} else {
		log.Println(successMsg)
	}
}

//定时器
func lngTimer() {
	config := getLogFromIni()
	timeFrom, err := strconv.Atoi(config.from)
	if err != nil {
		log.Fatalf("timeFrom类型转化错误,%s", err.Error())
	}
	timeTo, err := strconv.Atoi(config.timeConfig.to)
	if err != nil {
		log.Fatalf("timeTo类型转换错误,%s", err.Error())
	}
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			now := time.Now()
			fromTime := time.Date(now.Year(), now.Month(), now.Day(), timeFrom, 0, 0, 0, now.Location())
			toTime := time.Date(now.Year(), now.Month(), now.Day(), timeTo, 0, 0, 0, now.Location())
			if now.After(fromTime) && now.Before(toTime) {
				handleLngweatherLog()
			}
		}
	}
}

//发送邮件
func sendLngErrLog(username, password, host, to, body string) {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", username, password, hp[0])
	content_type := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + username + "\r\nSubject: lngweather发送失败\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, username, send_to, msg)
	if err != nil {
		log.Fatalf("邮件发送失败,%s", err.Error())
	}
}

func main() {
	lngTimer()
}
