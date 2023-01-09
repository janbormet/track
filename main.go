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
var storageType string

const fileTimew = "file-timew"
const fileJSON = "file-json"

func main() {

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startCmd.StringVar(&tags, "tags", "", "Supply tags as comma separated list. Example: a,b,c")
	startCmd.StringVar(&storageType, "storage", "file-json", "Supply the type of storage to be used (file-json or file-timew)")
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	stopCmd.StringVar(&storageType, "storage", "file-json", "Supply the type of storage to be used (file-json or file-timew)")

	_ = flag.NewFlagSet("stop", flag.ExitOnError)
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Use command `start` or `stop`\n")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "start":
		start(startCmd)
	case "stop":
		stop(stopCmd)
	}
}

func start(startCmd *flag.FlagSet) {
	_ = startCmd.Parse(os.Args[2:])
	switch storageType {
	case "file-json":
		err := tracker.Start(getContext(), types.MakeTags(strings.Split(tags, ",")...))
		if err != nil {
			panic(err)
		}

	}

}

func stop(stopCmd *flag.FlagSet) {
	_ = stopCmd.Parse(os.Args[2:])
	err := tracker.Stop(getContext())
	if err != nil {
		panic(err)
	}
}

func getContext() tracker.Context {
	switch storageType {
	case fileJSON:
		return tracker.Context{
			Storage: storage.NewDefaultJSONFileStorage(),
			Time:    time.Now(),
		}
	case fileTimew:
		return tracker.Context{
			Storage: storage.NewDefaultTimewFileStorage(),
			Time:    time.Now(),
		}
	default:
		return tracker.Context{
			Storage: storage.NewDefaultJSONFileStorage(),
			Time:    time.Now(),
		}
	}
}
