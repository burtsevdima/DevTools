package main

import (
	"fmt"
	"os"

	consts "github.com/Shohsta73/DevTools/Cher/constants"
)

func helpMessage() {
	fmt.Println("usage: Cher <command>\n" +
		"commands:\n" +
		"init | i   initialize Cher\n" +
		"help | h   this help message",
	)
}

func main() {
	CHER_DEBUG := os.Getenv("CHER_DEBUG")

	fmt.Println(consts.DEBUG)
	if CHER_DEBUG == "1" {
		consts.DEBUG = true
		fmt.Println(consts.DEBUG)
	}

	if consts.DEBUG {
		for i := 0; i < len(os.Args); i++ {
			fmt.Println("argument", i, ":", os.Args[i])
		}
		fmt.Println("len(os.Args):", len(os.Args))
		fmt.Println()
	}

	if len(os.Args) < 2 {
		helpMessage()
		return
	}

	switch os.Args[1] {
	case "help":
		fallthrough
	case "h":
		helpMessage()
		return
	case "init":
		fallthrough
	case "i":
		switch os.Args[2] {
		case "--help":
			fallthrough
		case "-h":
			fmt.Println(
				"usage: Cher init | i <flags> [langugaes]\n" +
					"flags:\n" +
					"--help | -h   this help message",
			)
			return
		}
	}
}
