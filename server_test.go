package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

func TestServe(t *testing.T) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	ports := make(chan HostPort, 1)
	go Serve(&wait, quit, ports)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	/* Wait until the program receives an interrupt */
	interrupt := <-sig
	log.Println(interrupt)
	/* Tell the servers to quit */
	log.Println("Quit the server threads")
	quit <- 1
	/* Wait for the servers to quit */
	log.Println("Wait for the server threads")
	wait.Wait()
	log.Println("Done waiting, Goodbye!")
}

func TestServeAndGoQuery(t *testing.T) {
	quit := make(chan int)
	wait := sync.WaitGroup{}
	ports := make(chan HostPort, 1)
	go Serve(&wait, quit, ports)
	/* channel to receive OS interrupts on */
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var hp HostPort
	hp = <- ports

	/* Now try to use goquery to scrape hosta */
	doc, err := goquery.NewDocument(fmt.Sprintf("http://%s:%d", hp.Host, hp.Port))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded doc ", doc.Url.Host)
	log.Println("Here's the text: ", doc.Text())

	/* Wait until the program receives an interrupt */
	interrupt := <-sig
	log.Println(interrupt)
	/* Tell the servers to quit */
	log.Println("Quit the server threads")
	quit <- 1
	/* Wait for the servers to quit */
	log.Println("Wait for the server threads")
	wait.Wait()
	log.Println("Done waiting, Goodbye!")
}
