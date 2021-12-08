package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	Host  *string
	Port  *string
	Rules *string
}

func ParseFlags() Flags {
	s := Flags{
		Rules: flag.String("rules", "", "File rule to be loaded"),
		Host:  flag.String("host", "127.0.0.1", "Host address to serve"),
		Port:  flag.String("port", "8000", "Host port to serve"),
	}

	flag.Parse()

	if len(*s.Rules) == 0 {
		fmt.Println("Use: -rules rule-file.json [OTHER_OPTIONAL_FLAGS]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return s
}
