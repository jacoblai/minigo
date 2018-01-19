package main

import (
	"runtime"
	"strconv"
	"fmt"
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"net/http"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		port = flag.Int("p", 9007, "web api port")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return
	}
	dir = dir + "/data"
	if _, err := os.Stat(dir); err != nil {
		log.Println(err)
		return
	}
	var y filecore.ConfigYml
	yml, _ := ioutil.ReadFile(dir + "/conf.json")
	cyml := filecore.MsgDecode(yml)
	err = json.Unmarshal(cyml, &y)
	if err != nil {
		log.Println(err)
	}
	filecore.AllowOrigins = y.AllowOrigin
	filecore.EnableAuth = y.VmAuth

	fdb := filecore.NewFileDb(*httpMode, *syncMode, *ipcPort, *limitupload)
	err = fdb.Open(y, dir)
	if err != nil {
		fdb.LogSend("yiyifd", "错误", err.Error())
		log.Println(err)
		return
	}

	router := httprouter.New()
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte("yiyifd service is on working..."))
	})
	router.GET("/api/dl/:fid", filecore.GzipHandler(fdb.PhotoGet))
	router.GET("/api/info/:fid", fdb.PhotoInfoGet)
	router.POST("/api/upload/:md5", fdb.Auth(fdb.PhotoPost))
	router.POST("/api/form", fdb.Auth(fdb.PhotoPostForm))

	srv := &http.Server{Handler: filecore.CORS(router)}
	if *tlsPort == 0 {
		srv.Addr = ":" + strconv.Itoa(*port)
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				fdb.LogSend("yiyifd", "错误", err.Error())
			}
		}()
		fdb.LogSend("yiyifd", "日志", "server on port "+strconv.Itoa(*port))
		fmt.Println("server on port", *port)
	} else {
		cert, err := tls.LoadX509KeyPair(dir+"/server.pem", dir+"/server.key")
		if err != nil {
			fdb.LogSend("yiyifd", "警告", err.Error())
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		srv.TLSConfig = config
		srv.Addr = ":" + strconv.Itoa(*tlsPort)
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				fdb.LogSend("yiyifd", "错误", err.Error())
			}
		}()
		fdb.LogSend("yiyifd", "日志", "server on tls port "+strconv.Itoa(*tlsPort))
		log.Println("server on tls port", *tlsPort)
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
			srv.Shutdown(ctx)
			fdb.Die()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
