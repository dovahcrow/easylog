package easylog

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"os"
	"strings"
)

var confName = `conf/cms.conf`
var log *logs.BeeLogger
var (
	LogCapacity    int64
	FileLog        bool
	FileSetting    string
	MailLog        bool
	MailSetting    string
	ConsoleLog     bool
	ConsoleSetting string
	ConnLog        bool
	ConnSetting    string
)

func init() {
	//defer func() { go log.StartLogger() }()
	conf, err := config.NewConfig(`ini`, confName)
	if err != nil {
		fmt.Println(`open config file fail:`, err)
		os.Exit(1)
	}

	// read&&set log capacity
	logcap, err := conf.Int64(`log::capacity`)
	if err != nil {
		fmt.Println(`Read log capacity error, use 1024 by default. Err:`, err)
		logcap = 1024
	}

	LogCapacity = logcap //export it

	log = logs.NewLogger(logcap)
	log.SetLevel(0)
	// read&&set log type
	// console log
	func() {
		ConsoleLog = false
		use, err := conf.Bool(`log-console::use`)
		if err != nil {
			fmt.Println(`Read console log section fail. Use console by default`)
			log.SetLogger(`console`, ``)
			return
		}
		if !use {
			return
		}
		ConsoleLog = true

		loglevel := getLogLevel(conf.String(`log-console::level`))

		ConsoleSetting = fmt.Sprintf("{\"level\":%d}", loglevel)
		log.SetLogger(`console`, ConsoleSetting)

	}()

	//file log
	func() {
		FileLog = false
		use, err := conf.Bool(`log-file::use`)
		if err != nil {
			fmt.Println(`Read file log section fail. Disable it`)
			return
		}
		if !use {

			return
		}
		FileLog = true

		//level
		loglevel := getLogLevel(conf.String(`log-file::level`))

		//filename
		logfile := conf.String(`log-file::filename`)
		if logfile == `` {
			fmt.Println(`Read filename fail. Use log.log by default`)
			logfile = `log.log`
		}

		//maxlines
		maxlines, err := conf.Int64(`log-file::maxlines`)
		if err != nil {
			fmt.Println(`read file log maxlines fail.use 1000000 by default`)
			maxlines = 1000000
		}

		//max file size
		maxsizes, err := conf.Int64(`log-file::maxsize`)
		if err != nil {
			fmt.Println(`Read file log maxsize fail. Use 256MB by default`)
			maxlines = 1 << 28
		}

		//rotate
		var rotate bool
		var dailyRotate bool
		var maxdays int64
		func() {
			rotate, err = conf.Bool(`log-file::rotate`)
			if err != nil {
				fmt.Println(`Read file rotate fail. Rotate off by default`)
				rotate = false
				return
			}
			if !rotate {
				return
			}
			dailyRotate, err = conf.Bool(`log-file::daily-rotate`)
			if err != nil {
				fmt.Println(`Read file daily-rotate fail. Rrotate by default`)
				dailyRotate = true
			}
			maxdays, err = conf.Int64(`log-file::maxdays`)
			if err != nil {
				fmt.Println(`Read file maxdays fail. Use 7 by default`)
				maxdays = 7
			}
		}()
		FileSetting = fmt.Sprintf(
			"{\"filename\":\"%s\",\"maxlines\":%d,\"maxsize\":%d,\"daily\":%v,\"maxdays\":%d,\"rotate\":%v,\"level\":%d}",
			logfile, maxlines, maxsizes, dailyRotate, maxdays, rotate, loglevel)
		log.SetLogger(`file`, FileSetting)
	}()
	func() {
		MailLog = false
		var use bool
		var username, password, host, sendTo, subject string
		var level int
		if use, err = conf.Bool(`log-mail::use`); err != nil {
			fmt.Println(`Read mail log error. Off by default`)
			return
		} else if !use {
			return
		}
		MailLog = true
		if username = conf.String(`log-mail::username`); username == `` {
			fmt.Println(`Read mail log username error. Off by default`)
			return
		}
		if password = conf.String(`log-mail::password`); password == `` {
			fmt.Println(`Read mail log password error. Off by default`)
			return
		}
		if host = conf.String(`log-mail::host`); host == `` {
			fmt.Println(`Read mail log host error.Off by default`)
			return
		}
		if sendTo = conf.String(`log-mail::sendTo`); sendTo == `` {
			fmt.Println(`Read mail log sendto error. Off by default`)
			return
		}
		if subject = conf.String(`log-mail::subject`); subject == `` {
			fmt.Println(`Read mail log subject error. Use "Diagnostic message from server" by default`)
		}
		level = getLogLevel(conf.String(`log-mail::level`))
		MailSetting = fmt.Sprintf(
			"{\"level\":%d,\"subject\":\"%s\",\"username\":\"%s\",\"password\":\"%s\",\"host\":\"%s\",\"sendTos\":[\"%s\"]}",
			level, subject, username, password, host, sendTo)
		log.SetLogger("smtp", MailSetting)
	}()

}

//re-export log functions
func Error(s string, v ...interface{})    { log.Error(s, v...) }
func Trace(s string, v ...interface{})    { log.Trace(s, v...) }
func Info(s string, v ...interface{})     { log.Info(s, v...) }
func Debug(s string, v ...interface{})    { log.Debug(s, v...) }
func Warn(s string, v ...interface{})     { log.Warn(s, v...) }
func Critical(s string, v ...interface{}) { log.Critical(s, v...) }

//re-export to an object(I need singleton!)
type l struct{}

var L = &l{}

func (*l) Error(s string, v ...interface{})    { Error(s, v...) }
func (*l) Trace(s string, v ...interface{})    { Trace(s, v...) }
func (*l) Info(s string, v ...interface{})     { Info(s, v...) }
func (*l) Debug(s string, v ...interface{})    { Debug(s, v...) }
func (*l) Warn(s string, v ...interface{})     { Warn(s, v...) }
func (*l) Critical(s string, v ...interface{}) { Critical(s, v...) }

func getLogLevel(loglevel string) int {
	//read&&set log level
	loglevel = strings.ToLower(loglevel)
	switch loglevel {
	case `trace`:
		{
			return 0
		}
	case `debug`:
		{
			return 1
		}
	case `info`:
		{
			return 2
		}
	case `warn`:
		{
			return 3
		}
	case `error`:
		{
			return 4
		}
	case `critical`:
		{
			return 5
		}
	default:
		{
			fmt.Println(`Read log level fail. Use trace by default`)
			return 0
		}
	}
}
