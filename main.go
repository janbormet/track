package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"track/storage"
	"track/tracker"
	"track/types"
)

var tags string

func main() {

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startCmd.StringVar(&tags, "tags", "", "Supply tags as comma separated list. Example: a,b,c")

	_ = flag.NewFlagSet("stop", flag.ExitOnError)
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Use command `start` or `stop`\n")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "start":
		start(startCmd)
	case "stop":
		stop()
	}
}

func start(startCmd *flag.FlagSet) {
	_ = startCmd.Parse(os.Args[2:])
	fmt.Println(strings.Split(tags, ","))
	err := tracker.Start(defaultContext(), types.MakeTags(strings.Split(tags, ",")...))
	if err != nil {
		panic(err)
	}
}

func stop() {
	err := tracker.Stop(defaultContext())
	if err != nil {
		panic(err)
	}
}

func defaultContext() tracker.Context {
	return tracker.Context{
		Storage: storage.NewDefaultFileStorage(),
		Time:    time.Now(),
	}
}
