package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/braintree/manners"
)

func init() {
	http.HandleFunc("/contributors/", handler)
}

func Run() {
	host := "0.0.0.0:5000"
	s := manners.NewWithServer(&http.Server{
		Addr:           host,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	})

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan

		s.Close()

		go func() {
			time.Sleep(time.Second * 60)
			fmt.Println("Shutting down server but not gracefully")
			os.Exit(1)
		}()
	}()

	log.Println("Waiting for connections on: ", host)
	s.ListenAndServe()
}
