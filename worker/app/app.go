package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"service/pkg/config"
	"service/pkg/db"
	"service/pkg/log"
	"service/worker/components"
	"strings"
	"syscall"
	"time"
)

//App ..
type App struct {
	logFile  string
	confFile string
	PFile    *os.File
}

//New ...
func New(logFile string, confFile string, file *os.File) *App {
	return &App{
		logFile:  logFile,
		confFile: confFile,
		PFile:    file,
	}
}

//Run ...
func (app *App) Run() {
	//Init config
	config.Init(app.confFile)

	//Initalize logger
	{
		app.logFile = strings.Replace(app.logFile, "%s", time.Now().Format("20060102150405"), 1)
		log.Init(app.logFile, config.GetLogLevel())
	}

	db.SetConfigDB(config.LoadDBConf("DB"))
	db.SetConfigDBS(config.LoadDBConf("DBS"))
	db.Init()
	db.InitS()

	name := flag.String("name", "Service Worker", "name to print")
	flag.Parse()
	log.Printf("Starting sleepservice for %s", *name)
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	// signal.Notify(sigs)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	//signal.Notify(sigs,syscall.SIGQUIT)
	// method invoked upon seeing signal

	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		app.PFile.Close()

		err := os.Remove(app.PFile.Name())
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	}()

	Nsecs := 3000
	f := true
	t := true

	for {

		//sleep Nsecs
		time.Sleep(time.Millisecond * time.Duration(Nsecs))

		timeStamp := time.Now()
		today := timeStamp.Day()

		timestr := timeStamp.Format("15:04")

		//
		components.UpdateStatusSubscription()

		//Every day in TIME_DAILY_ONE and TIME_DAILY_TWO send sms status if client have violation
		{
			if timestr == config.GetString("TIME_DAILY_ONE") || timestr == config.GetString("TIME_DAILY_TWO") {
				components.InsertDailySubscription()
			}
		}

		//Monthly every n-th day send sms in time MONTHLY_NTH_TIME
		{
			if today%config.GetInt("MONTHLY_NTH_DAY") == 0 && timestr == config.GetString("MONTHLY_NTH_TIME") && f {
				components.InsertMonthlySubscription()
				f = false
			}

			if today%config.GetInt("MONTHLY_NTH_DAY") != 0 || timestr != config.GetString("MONTHLY_NTH_TIME") {
				f = true
			}
		}

		//Every day in END_SUBSCRIPTION_TIME sent sms notification about ending subscription
		{
			if timestr == config.GetString("END_SUBSCRIPTION_TIME") && t {
				components.InsertEndingSubscription(config.GetInt("DAYS_AGO"))
				t = false
			}

			if timestr != config.GetString("END_SUBSCRIPTION_TIME") {
				t = true
			}
		}

		//Every Nsecs proceed queue
		components.ProceedQueue(context.Background())

	}

}
