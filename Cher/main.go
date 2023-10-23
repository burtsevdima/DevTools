package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	consts "github.com/Shohsta73/DevTools/Cher/constants"
	"github.com/Shohsta73/DevTools/Cher/parser"
)

func helpMessage() {
	fmt.Println("usage: Cher <command>\n" +
		"commands:\n" +
		"init | i   initialize Cher\n" +
		"help | h   this help message",
	)
}

func initCher() (string, error) {
	configDir, err := getConfigDir()
	if configDir == "" {
		return "", errors.New("config directory not found")
	}
	if err != nil {
		return "", err
	}

	_, err = os.Stat(configDir)

	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err = os.Mkdir(configDir, 0755) // 0755 is the directory permission, you can adjust it as needed
		if err != nil {
			return "", err
		}

		fmt.Printf(
			"We do not provide default .editorconfig files.\n"+
				"You will have to create your own in %s.\n"+
				"We recomande ussing directories <lang> for specif languages.\n", configDir,
		)
	} else if err != nil {
		return "", err // Some other error occurred
	}

	return configDir, nil
}

func getConfigDir() (string, error) {
	var configDir string
	userDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if runtime.GOOS == "windows" {
		configDir = userDir + "\\Cher"
	} else {
		configDir = userDir + "/Cher"
	}

	return configDir, nil
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

	Parser := parser.NewParser()

	configDir, err := getConfigDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(configDir)

	parsedCommands, err := Parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	if parsedCommands.ParsedCommands["help"] {
		helpMessage()
		return
	}

	if parsedCommands.ParsedCommands["init"] {
		_, err := initCher()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if parsedCommands.ParsedCommands["add"] {
		var DirSep string
		if runtime.GOOS == "windows" {
			DirSep = "\\"
		} else {
			DirSep = "/"
		}
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "help" {
				helpMessage()
				return
			}

			langDir := configDir + DirSep + os.Args[i]
			_, err = os.Stat(langDir)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		return
	}

	if parsedCommands.ParsedCommands["new"] {
		var DirSep string
		if runtime.GOOS == "windows" {
			DirSep = "\\"
		} else {
			DirSep = "/"
		}
		return
	}

	if parsedCommands.ParsedCommands["remove"] {
		var DirSep string
		if runtime.GOOS == "windows" {
			DirSep = "\\"
		} else {
			DirSep = "/"
		}
		return
	}

	if parsedCommands.ParsedCommands["edit"] {
		var DirSep string
		if runtime.GOOS == "windows" {
			DirSep = "\\"
		} else {
			DirSep = "/"
		}
		return
	}
}
