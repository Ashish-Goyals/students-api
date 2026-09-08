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
)




func main(){
	// fmt.Println("Welcome to students api")
	// load config 
        cfg :=config.MustLoad()

		// databse setup 

	// setup router 
       router := http.NewServeMux()
	   router.HandleFunc("GET /", func(w http.ResponseWriter,r *http.Request){
		w.Write([]byte("Welcome to students api"))
	   })
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

	 slog.Info("shuutting down the server")

	 ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	 defer cancel()


	 err := server.Shutdown(ctx)

	 if err != nil {
		slog.Error("Failed to Shutdown srver", slog.String("error",err.Error()))
	 }
     slog.Info("server sgutdoen succesfully")
	

}
