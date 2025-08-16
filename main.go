package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const DonVersion = "0.0.1"

type Command struct {
	Name string
	Desc string
}

var commands = []Command{
	{"com", "Commit all changes with timestep and push to currect branch"},
	{"tail", "Run TailwindCSS watcher"},
}

var donCommands = []Command{
	{"version", "Print the currect don version"},
	{"update", "Update don to the latest version from master"},
}

func stdInsert(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func xrr(msg string, err error) {
	if err != nil {	
		fmt.Println(msg, err)
	} else {
		fmt.Println(msg)
	}
	os.Exit(2)
}

func printHelp() {
	fmt.Println("Usage: don <command> [arguments]")
	fmt.Println("\ndon commands:")
	for _, cmd := range donCommands {
		fmt.Printf("  %-8s %s\n", cmd.Name, cmd.Desc)
	}
	fmt.Println("\nCommands:")
	for _, cmd := range commands {
		fmt.Printf("  %-8s %s\n", cmd.Name, cmd.Desc)
	}
}

func main() {
	var err error

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(2)
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("don version: %s\n", DonVersion)

	case "update":
		cmdUpdate := exec.Command("go", "install", "github.com/Zimzozaur/don@master")
		stdInsert(cmdUpdate)
		if err = cmdUpdate.Run(); err != nil {
			xrr("go install failed: %v", err)
		}

	case "com":
		cmdAdd := exec.Command("git", "add", ".")
		stdInsert(cmdAdd)
		if err = cmdAdd.Run(); err != nil {
			xrr("git add failed: %v", err)
			
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		cmdCommit := exec.Command("git", "commit", "-m", "don: "+timestamp)
		stdInsert(cmdCommit)
		if err = cmdCommit.Run(); err != nil {
			xrr("git commit failed: %v", err)
		}

		out, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		branch := strings.TrimSpace(string(out))

		cmdPush := exec.Command("git", "push", "origin", branch)
		stdInsert(cmdPush)
		if err = cmdPush.Run(); err != nil {
			xrr("git push failed: %v", err)
		}

	case "tail":
		if len(os.Args) < 3 {
			 xrr("Usage: don tail <path>", nil)
		}
		tailPath := os.Args[2]
		cmdTail := exec.Command("tailwindcss", "-i", tailPath+"/input.css", "-o", tailPath+"/output.css", "--watch")
		stdInsert(cmdTail)
		if err = cmdTail.Run(); err != nil {
			xrr("tailwindcss failed: %v", err)
		}
		
	default:
		fmt.Printf("don %s: unknown command\n", os.Args[1])
		fmt.Println("Run 'don' to list commands.")
		os.Exit(2)
	}
}
