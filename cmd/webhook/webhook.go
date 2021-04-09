package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	model "github.com/etng/prom_alertmanager_webhook_receiver/model"
	"github.com/etng/prom_alertmanager_webhook_receiver/notifier"
	"github.com/gin-gonic/gin"
)

var (
	showHelp              bool
	defaultDingtalkToken  string
	defaultDingtalkPrefix string
	listenPort            int
	debug                 bool
	logFilename           string
)
var options *notifier.NotifyOptions
var dingding *notifier.DingTalkNotifier
var pushPlus *notifier.PushPlusNotifier

func parseOptions() {
	options = &notifier.NotifyOptions{}
	flag.BoolVar(&showHelp, "h", false, "help")
	flag.StringVar(&logFilename, "log", "", "write log to this filename")
	flag.StringVar(&options.DingTalkToken, "token", "", "default dingding robot token")
	flag.StringVar(&options.DingTalkPrefix, "prefix", "", "default dingding robot content prefix")
	flag.StringVar(&options.PushPlusToken, "pp_token", "", "default push plus token")
	flag.StringVar(&options.PushPlusGroup, "pp_topic", "", "default push plus topic")
	flag.BoolVar(&debug, "debug", false, "trun debug on")
	flag.IntVar(&listenPort, "port", 8080, "port to listen")
	flag.Parse()
}

func main() {
	parseOptions()
	if showHelp {
		flag.Usage()
		return
	}
	dingding = notifier.NewDingTalkNotifier(options.DingTalkToken, options.DingTalkPrefix)
	pushPlus = notifier.NewPushPlusNotifier(options.PushPlusToken, options.PushPlusGroup)
	var logWriter io.Writer
	if logFilename != "" {
		os.MkdirAll(filepath.Dir(logFilename), 0777)
		if of, e := os.OpenFile(logFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); e == nil {
			if debug {
				logWriter = io.MultiWriter(os.Stdout, of)
			} else {
				logWriter = of
			}
		} else {
			log.Printf("fail to open log file %s for %s", logFilename, e)
		}
		if logWriter != nil {
			log.SetOutput(logWriter)
			gin.DefaultWriter = logWriter
		}
	}

	gin_mode := gin.ReleaseMode
	if debug {
		gin_mode = gin.DebugMode
	}
	gin.SetMode(gin_mode)
	router := gin.Default()
	router.POST("/", ReceiveHandler)
	router.POST("/webhook", ReceiveHandler)
	router.Run(fmt.Sprintf(":%d", listenPort))
}

func ReceiveHandler(c *gin.Context) {
	if debug {
		req, _ := httputil.DumpRequest(c.Request, true)
		log.Printf("received alert \n%s", req)
	}
	var notification model.Notification

	err := c.BindJSON(&notification)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	errMsg := ""
	err = dingding.Send(notification)
	if err != nil {
		log.Printf("fail to send dingtalk notification for %s", err)
		errMsg = fmt.Sprintf("%s;fail to send dingtalk notification for %s", errMsg, err)
	}
	err = pushPlus.Send(notification)

	if err != nil {
		log.Printf("fail to send pushPlus notification for %s", err)
		errMsg = fmt.Sprintf("%s;fail to send pushPlus notification for %s", errMsg, err)
	}
	errMsg = strings.TrimPrefix(errMsg, ";")

	if errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("notification fail: %s", errMsg),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "send notifcation successfully!"})
}
