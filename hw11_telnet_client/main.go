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

	if err := telnet.Connect(); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	telnet.Start(ctx, stop)

	<-ctx.Done()
}
