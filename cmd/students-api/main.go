package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashish-goyals/students-api/internal/config"
	"github.com/ashish-goyals/students-api/internal/http/handlers/student"
	"github.com/ashish-goyals/students-api/internal/storage/sqlite"
)




func main(){
	// fmt.Println("Welcome to students api")
	// load config 
        cfg :=config.MustLoad()

		// databse setup 
storage, err := sqlite.New(cfg)

		if err != nil {
			log.Fatal(err)
		}
		slog.Info("Storage Intialised",slog.String("env",cfg.Env),slog.String("version","1.0.0"))	

	// setup router 
       router := http.NewServeMux()
	   router.HandleFunc("POST /api/students",student.New(storage))
	// setup server 
       server :=  http.Server{
		Addr: cfg.Addr,
		Handler: router,

	  }
	  slog.Info("server started ", slog.String("address",cfg.HTTPServer.Addr))
	    fmt.Printf("Server started %s ", cfg.HTTPServer.Addr)

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT,syscall.SIGTERM)


     go func(){
            err := server.ListenAndServe()
	       if err != nil{
		   log.Fatal("Failed to start server ",err)
	    } 
       }()
	
     <-done

	 slog.Info("shutting down the server")

	 ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	 defer cancel()


	 err = server.Shutdown(ctx)

	 if err != nil {
		slog.Error("Failed to Shutdown srver", slog.String("error",err.Error()))
	 }
     slog.Info("server sgutdoen succesfully")
	

}
