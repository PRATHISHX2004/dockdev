package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"generator/internal"
)

func main() {
	args := os.Args

	if len(args) == 3 && args[1] == "rm" {
		internal.DeleteProject(args[2])
		return
	}

	var domain string
	if len(args) == 2 {
		domain = args[1]
	} else {
		fmt.Print("Enter project domain (e.g. app.test): ")
		reader := bufio.NewReader(os.Stdin)
		d, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read input:", err)
			os.Exit(1)
		}
		domain = strings.TrimSpace(d)
	}

	if domain == "" {
		fmt.Println("Domain cannot be empty.")
		os.Exit(1)
	}

	if err := internal.GenerateProject(domain); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}