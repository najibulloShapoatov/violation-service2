package app

import (
	"context"
	"service/pkg/config"
	"service/pkg/db"
	"service/pkg/log"
	"service/web/router"

	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/cors"
)

//App ...
type App struct {
	logFile  string
	confFile string
	HTTPPORT string
	PFile    *os.File
}

//NewApp ..
func NewApp(logFile, confFile string, file *os.File) *App {
	return &App{
		logFile:  logFile,
		confFile: confFile,
		HTTPPORT: "-",
		PFile:    file,
	}
}

//Run ...
func (a *App) Run() {

	//Init config
	config.Init(a.confFile)

	//Initalize logger
	{
		a.logFile = strings.Replace(a.logFile, "%s", time.Now().Format("20060102150405"), 1)
		log.Init(a.logFile, config.GetLogLevel())
	}

	db.SetConfigDB(config.LoadDBConf("DB"))
	db.SetConfigDBS(config.LoadDBConf("DBS"))
	db.Init()
	db.InitS()

	a.HTTPPORT = config.GetHTTPPort()

	var router = router.Init()

	//creaet http server
	srv := &http.Server{
		Addr:         ":" + a.HTTPPORT,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      cors.AllowAll().Handler(router),
	}

	log.Info("Listining on: ", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Warn("Closing runned server ", err, fmt.Sprintf("%+v", srv))
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.

	a.PFile.Close()

	err := os.Remove(a.PFile.Name())
	if err != nil {
		fmt.Println(err)
	}

	log.Info("\n <----\tshutting down App\t---->\n")
	os.Exit(0)

}
