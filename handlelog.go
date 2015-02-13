package main

import (
	"bufio"
	"github.com/Unknwon/goconfig"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

//读取配置文件
func getLogFromIni() (string, string, string, string, string) {
	configFile := "config.ini"
	key := "logPath"
	section := "mail"
	config, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Fatalf("载入配置文件%s错误", configFile)
	}
	logPath, err := config.GetValue(goconfig.DEFAULT_SECTION, key)
	if err != nil {
		log.Fatalf("读取节点%s错误", key)
	}
	username, err := config.GetValue(section, "username")
	if err != nil {
		log.Fatalf("读取节点%s错误", "username")
	}
	password, err := config.GetValue(section, "password")
	if err != nil {
		log.Fatalf("读取节点%s错误", "password")
	}
	host, err := config.GetValue(section, "host")
	if err != nil {
		log.Fatalf("读取节点%s错误", "host")
	}
	to, err := config.GetValue(section, "to")
	if err != nil {
		log.Fatalf("读取节点%s错误", "to")
	}
	return logPath, username, password, host, to
}

//处理lngweather日志
func handleLngweatherLog() {
	fileName, username, password, host, to := getLogFromIni()
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalf("读入文件%s出错", fileName)
	}
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
		sendLngErrLog(username, password, host, to, errMsg)
	} else {
		log.Println(successMsg)
	}
}

//定时器
func lngTimer() {
	timer := time.NewTicker(10 * time.Hour)
	for {
		select {
		case <-timer.C:
			now := time.Now()
			beforeTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location())
			afterTime := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())
			if now.After(beforeTime) && now.Before(afterTime) {
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
	smtp.SendMail(host, auth, username, send_to, msg)
}

func main() {
	lngTimer()
}
