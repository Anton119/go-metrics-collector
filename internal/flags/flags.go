package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	FlagRunAddr        string
	FlagPollInterval   int
	FlagReportInterval int
)

func ParseFlags() {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port of HTTP server")
	fs.IntVar(&FlagPollInterval, "p", 2, "poll interval in seconds")
	fs.IntVar(&FlagReportInterval, "r", 10, "report interval in seconds")

	// ищем неизвестные флаги
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// проверка на неизвестные флаги
	if fs.NArg() > 0 {
		fmt.Printf("Unknown arguments: %v\n", fs.Args())
		os.Exit(1)
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if IntEnvReportInterval, err := strconv.Atoi(envReportInterval); err == nil && IntEnvReportInterval > 0 {
			FlagReportInterval = IntEnvReportInterval
		} else if err != nil {
			log.Printf("invalid REPORT_INTERVAL: %v", err)
		}

	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if IntEnvPollInterval, err := strconv.Atoi(envPollInterval); err == nil && IntEnvPollInterval > 0 {
			FlagPollInterval = IntEnvPollInterval
		} else if err != nil {
			log.Printf("invalid POLL_INTERVAL: %v", err)
		}
	}

}
