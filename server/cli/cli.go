package main

import (
	"fmt"
	"os"
	"server/cli/export"
	"server/cli/load"
	"server/cli/prepare"
	_ "server/services/database"
)

const cliName = "ojtool"

func help() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Printf("  %s <cmd>\n", cliName)
	fmt.Println("valid cmd:")
	fmt.Println("  prepare: prepare Tp0tOJ environment.")
	fmt.Println("  export: export OJ data to file.")
	fmt.Println("  load: auto register users from file.")
	fmt.Println("  help: show help information.")
	fmt.Printf("Use `%s help <cmd>` for detailed help information.\n", cliName)
}

func main() {
	cmdList := map[string]func(args []string){
		"prepare": prepare.Run,
		"export":  export.Run,
		"load":    load.Run,
	}

	if len(os.Args) == 1 {
		help()
		return
	}
	if os.Args[1] == "help" {
		if len(os.Args) == 3 {
			if fn, ok := cmdList[os.Args[2]]; ok {
				fn([]string{"-help"})
				return
			} else {
				fmt.Printf("\033[31;1mERROR: unknown cmd \"%s\"\033[0m\n", os.Args[1])
			}
		} else if len(os.Args) != 2 {
			fmt.Println("\033[31;1mERROR: wrong number of arguments\033[0m")
		}
		help()
		return
	}
	if fn, ok := cmdList[os.Args[1]]; ok {
		fn(os.Args[2:])
	} else {
		fmt.Printf("\033[31;1mERROR: unknown cmd \"%s\"\033[0m\n", os.Args[1])
		help()
	}
}
