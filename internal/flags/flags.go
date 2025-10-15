package flags

import (
	"flag"
	"fmt"
	"os"
)

var (
	FlagRunAddr        string
	FlagPollInterval   int
	FlagReportInterval int
)

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// сервер
	flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
	// агент
	flag.IntVar(&FlagPollInterval, "p", 2, "poll interval in seconds")
	flag.IntVar(&FlagReportInterval, "r", 10, "report interval in seconds")

	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
