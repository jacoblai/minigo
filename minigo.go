package main

import (
	"runtime"
	"strconv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"
	"path/filepath"
	"io/ioutil"
	"net/http"
	"github.com/ghodss/yaml"
	"github.com/jacoblai/httprouter"
	"context"
	"engine"
	"auth"
	"cors"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		conf = flag.String("c", "/conf.yaml", "config yaml file flag")
	)
	flag.Parse()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = dir + "/data"
	if _, err := os.Stat(dir); err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	ymlfile := ""
	if *conf == "/conf.yaml" {
		ymlfile = dir + *conf
	} else {
		ymlfile = *conf
	}

	//启动文件日志
	logFile, logErr := os.OpenFile(dir+"/dal.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		log.Printf("err: %v\n", logErr)
		fmt.Printf("err: %v\n", logErr)
		return

	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if _, err := os.Stat(ymlfile); err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}
	buf, _ := ioutil.ReadFile(ymlfile)
	var y engine.ConfigYml
	err = yaml.Unmarshal(buf, &y)
	if err != nil {
		log.Printf("err: %v\n", err)
		fmt.Printf("err: %v\n", err)
		return
	}

	eng := &engine.DbEngine{}
	eng.SystemConfig = &y
	err = eng.Open(dir)
	if err != nil{
		log.Fatal("database connect error")
	}

	router := httprouter.New()
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte("service is on working..."))
	})
	router.POST("/api/users", auth.Auth(eng.AddUser))

	srv := &http.Server{Handler: cors.CORS(router)}
	srv.Addr = ":" + strconv.Itoa(y.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("server on port", y.Port)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
			srv.Shutdown(ctx)
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
