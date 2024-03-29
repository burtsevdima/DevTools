package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	consts "github.com/Shohsta73/DevTools/Cher/constants"
	"github.com/Shohsta73/DevTools/Cher/parser"
	ini "gopkg.in/ini.v1"
)

func fallbackEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		switch runtime.GOOS {
		case "windows":
			editor = "notepad.exe"
		case "linux":
			editor = "nano"
		default:
			fmt.Printf("Unsupported OS: %s.\n"+
				"Please checkout the issuse here:", runtime.GOOS)
		}
	}

	return editor
}

func helpMessage() {
	fmt.Println("usage: Cher <command>\n" +
		"commands:\n" +
		"init | i   initialize Cher\n" +
		"help | h   this help message",
	)
}

func writeConfigFile(dirPath *string) {
	file := filepath.Join(*dirPath, "cher.config.ini")

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		// File doesn't exist, create it
		file, err := os.Create(file)
		if err != nil {
			fmt.Printf("Failed to create configuration file: %v\n", err)
		}
		file.Close() // Close the file to ensure it's created and empty
	}

	fmt.Print("Please enter the command for oppening text editor of your choice: ")
	var userInput string
	_, err = fmt.Scanln(&userInput)
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	// Ensure the user input is a valid index
	if userInput == "" {
		fmt.Printf("Invalid choice. Please input command for oppening text editro of your choice\n")
		return
	}

	fileOpened, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Failed to open configuration file: %v\n", err)
	}

	contents := "[editor]\n" + "command = " + userInput + "\n"

	_, err = fileOpened.WriteString(contents)
	if err != nil {
		fmt.Printf("Failed to write configuration file: %v\n", err)
	}
	fileOpened.Close()
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

	writeConfigFile(&configDir)

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
	if consts.DEBUG {
		fmt.Println(configDir)
	}

	configFile := filepath.Join(configDir, "cher.config.ini")

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

		if len(os.Args) < 3 {
			fmt.Println("Usage: add <lang(s)>")
			return
		}

		var contents string

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

			files, err := os.ReadDir(langDir)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Files in", langDir, "directory:")
			for i := 0; i < len(files); i++ {
				fileName := files[i].Name()
				extension := filepath.Ext(fileName)
				fileNameWithoutExt := fileName[:len(fileName)-len(extension)]
				fmt.Printf("%d. %s\n", (i + 1), fileNameWithoutExt)
			}

			fmt.Print("Enter the number of the file you want to edit: ")
			var userInput int
			_, err = fmt.Scanln(&userInput)
			if err != nil {
				fmt.Println("Invalid input:", err)
				return
			}

			// Ensure the user input is a valid index
			if userInput < 1 || userInput > len(files) {
				fmt.Println("Invalid choice. Please select a number between 1 and", len(files))
				return
			}

			// Delete the selected file
			selectedFile := files[userInput-1]
			content, err := os.ReadFile(selectedFile.Name())
			if err != nil {
				fmt.Println(err)
				continue
			}

			contents += string(content)
		}

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = os.WriteFile(filepath.Join(wd, ".editorconfig"), []byte(contents), 0644)
		if err != nil {
			fmt.Println("Error:", err)
			return
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

		if len(os.Args) < 3 {
			fmt.Println("Usage: new <lang(s)>")
			return
		}

		cfg, err := ini.Load(configFile)
		if consts.DEBUG {
			fmt.Println(cfg)
		}

		var editor string

		if err == nil {
			section, err := cfg.GetSection("editor")
			if err != nil {
				fmt.Println(err)
				editor = fallbackEditor()
			}

			command, err := section.GetKey("command")
			if err != nil {
				fmt.Println(err)
				editor = fallbackEditor()
			}
			editor = command.String()
		} else {
			editor = fallbackEditor()
		}

		for i := 2; i < len(os.Args); i++ {
			langDir := configDir + DirSep + os.Args[i]
			_, err := os.Stat(langDir)

			if os.IsNotExist(err) {
				fmt.Println("Directory", langDir, "does not exist.")
				continue
			} else {
				fmt.Print("Please enter the name for new language configuration file: ")
				var userInput string
				_, err = fmt.Scanln(&userInput)
				if err != nil {
					fmt.Println("Invalid input:", err)
					continue
				}

				if userInput == "" {
					fmt.Printf("Not creating new configuration file for %s as ther was no name provided\n", os.Args[i])
					continue
				}

				cmd := exec.Command(editor, langDir+string(os.PathSeparator)+userInput+".editorconfig")

				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Start()
				if err != nil {
					panic(err)
				}

				err = cmd.Wait()
				if err != nil {
					panic(err)
				}
			}
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

		if len(os.Args) < 3 {
			fmt.Println("Usage: remove <lang>")
			return
		}
		langDir := configDir + DirSep + os.Args[2]

		configs, err := os.ReadDir(langDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Files in", langDir, "directory:")
		for i := 0; i < len(configs); i++ {
			fileName := configs[i].Name()
			extension := filepath.Ext(fileName)
			fileNameWithoutExt := fileName[:len(fileName)-len(extension)]
			fmt.Printf("%d. %s\n", (i + 1), fileNameWithoutExt)
		}

		// Prompt the user to select a file
		fmt.Print("Enter the number of the file you want to delete: ")
		var userInput int
		_, err = fmt.Scanln(&userInput)
		if err != nil {
			fmt.Println("Invalid input:", err)
			return
		}

		// Ensure the user input is a valid index
		if userInput < 1 || userInput > len(configs) {
			fmt.Println("Invalid choice. Please select a number between 1 and", len(configs))
			return
		}

		// Delete the selected file
		selectedFile := configs[userInput-1]
		err = os.Remove(langDir + string(os.PathSeparator) + selectedFile.Name())
		if err != nil {
			fmt.Println("Error deleting file:", err)
		} else {
			fmt.Println("File", selectedFile.Name(), "has been deleted.")
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

		if len(os.Args) < 3 {
			fmt.Println("Usage: edit <lang>")
			return
		}

		cfg, err := ini.Load(configFile)
		if consts.DEBUG {
			fmt.Println(cfg)
		}

		var editor string

		if err == nil {
			section, err := cfg.GetSection("editor")
			if err != nil {
				fmt.Println(err)
				editor = fallbackEditor()
			}

			command, err := section.GetKey("command")
			if err != nil {
				fmt.Println(err)
				editor = fallbackEditor()
			}
			editor = command.String()
		} else {
			editor = fallbackEditor()
		}
		langDir := configDir + DirSep + os.Args[2]

		if consts.DEBUG {
			fmt.Println(langDir)
		}

		files, err := os.ReadDir(langDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Files in", langDir, "directory:")
		for i := 0; i < len(files); i++ {
			fileName := files[i].Name()
			extension := filepath.Ext(fileName)
			fileNameWithoutExt := fileName[:len(fileName)-len(extension)]
			fmt.Printf("%d. %s\n", (i + 1), fileNameWithoutExt)
		}

		fmt.Print("Enter the number of the file you want to edit: ")
		var userInput int
		_, err = fmt.Scanln(&userInput)
		if err != nil {
			fmt.Println("Invalid input:", err)
			return
		}

		// Ensure the user input is a valid index
		if userInput < 1 || userInput > len(files) {
			fmt.Println("Invalid choice. Please select a number between 1 and", len(files))
			return
		}

		// Delete the selected file
		selectedFile := files[userInput-1]

		cmd := exec.Command(editor, langDir+string(os.PathSeparator)+selectedFile.Name())

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		err = cmd.Wait()
		if err != nil {
			panic(err)
		}

		return
	}

	if parsedCommands.ParsedCommands["list"] {
		var DirSep string
		if runtime.GOOS == "windows" {
			DirSep = "\\"
		} else {
			DirSep = "/"
		}
		fmt.Printf("DirSep: %v\n", DirSep)

		langs, err := os.ReadDir(configDir)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < len(langs); i++ {
			fmt.Println(langs[i].Name())
		}
		return
	}
}
