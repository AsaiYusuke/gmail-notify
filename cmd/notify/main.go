package main

import (
	"log"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type config struct {
	LogMaxSize      int
	LogBackups      int
	LogMaxAge       int
	LogCompress     bool
	Interval        time.Duration
	RetryInterval   time.Duration
	AppID           string
	Title           string
	Labels          []string
	NotifyBeginHour int
	NotifyEndHour   int
}

func main() {
	p := paths{}
	p.init()

	conf := getConfig(p)

	log.SetOutput(&lumberjack.Logger{
		Filename:   p.getLogFilePath(),
		MaxSize:    conf.LogMaxSize,
		MaxBackups: conf.LogBackups,
		MaxAge:     conf.LogMaxAge,
		Compress:   conf.LogCompress,
	})

	terminateDuplicateProcess()

	shell := powerShell{}
	shell.open()

	template := commandTemplate{}
	template.init()

	statusList := make([]status, 0)

	for _, subDir := range p.getServiceDirs() {
		status := status{
			paths:  &p,
			conf:   &conf,
			subDir: subDir,
			idsFile: &jsonAccess{
				filename: p.getUnreadJSONFilePath(subDir),
			},
			gmail: &gmailAPI{
				paths:  &p,
				conf:   &conf,
				subDir: subDir,
			},
			toast: &toastNotify{
				shell:    &shell,
				template: &template,
			},
		}

		statusList = append(statusList, status)
	}

	for {
		nowHour := time.Now().Hour()
		if checkValidHour(nowHour, conf) || checkValidHour(nowHour+24, conf) {
			for _, status := range statusList {
				status.update()
			}
			time.Sleep(time.Minute * conf.Interval)
		} else {
			time.Sleep(getDurationToNextStart(conf))
		}
	}
}

func getConfig(p paths) config {
	conf := config{
		LogMaxSize:      50,
		LogBackups:      3,
		LogMaxAge:       28,
		LogCompress:     false,
		Interval:        5,
		RetryInterval:   60,
		AppID:           "Gmail notify",
		Title:           "You've got mail.",
		Labels:          []string{"INBOX", "UNREAD"},
		NotifyBeginHour: 7,
		NotifyEndHour:   22,
	}

	configJSON := jsonAccess{filename: p.getConfigJSONFilePath()}
	configJSON.readConfig(&conf)
	conf.NotifyBeginHour = conf.NotifyBeginHour % 24
	conf.NotifyEndHour = conf.NotifyEndHour % 24
	if conf.NotifyBeginHour > conf.NotifyEndHour {
		conf.NotifyEndHour += 24
	}
	return conf
}

func checkValidHour(hour int, conf config) bool {
	return hour >= conf.NotifyBeginHour && hour <= conf.NotifyEndHour
}

func getDurationToNextStart(conf config) time.Duration {
	hour, min, sec := time.Now().Clock()

	var offsetHour int
	if hour >= conf.NotifyBeginHour {
		offsetHour = 24
	}

	return time.Hour*time.Duration(offsetHour+conf.NotifyBeginHour-hour) -
		(time.Minute*time.Duration(min) + time.Second*time.Duration(sec))
}
