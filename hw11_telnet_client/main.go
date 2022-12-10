package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "server connection timeout")
	flag.Parse()

	telnet := NewTelnetClient(fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1)), *timeout, os.Stdin, os.Stdout)
	defer func() {
		if err := telnet.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := telnet.Connect(); err != nil {
		log.Println(err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	go func() {
		defer stop()
		if err := telnet.Send(); err != nil {
			log.Println(err)
		}
	}()
	go func() {
		defer stop()
		if err := telnet.Receive(); err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Bye-bye")
	}()

	<-ctx.Done()
}
