package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyokomi/emoji"

	"github.com/kgantsov/uptime/app/monitor"
)

func main() {
	configPathPtr := flag.String("config", "./config.json", "A path to a config.json file")

	flag.Parse()

	config, err := monitor.ReadConfig(*configPathPtr)
	if err != nil {
		fmt.Printf("Got an error parsing config %s", err)
	}

	dispatcher := monitor.NewDispatcher(config.Services)
	dispatcher.Start()

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("Got signal: %s\n", sig)

		fmt.Print("Stopping monitoring\n")

		dispatcher.Stop()

		time.Sleep(100 * time.Millisecond)

		done <- struct{}{}
	}()

	emoji.Printf("Started uptime monitor\n")

	<-done

	fmt.Print("Stopped monitoring\n")
}
