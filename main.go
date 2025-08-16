package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const DonVersion= "0.0.12"

func stdInsert(cmd *exec.Cmd)  {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func main() {
	var err error
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: don <command>")
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("don version: %s\n", DonVersion)
		
	case "update":
		cmdUpdate := exec.Command("go", "install", "github.com/Zimzozaur/don@main")
		cmdUpdate.Env = append(os.Environ(), "GOPROXY=direct")
		stdInsert(cmdUpdate)
		if err = cmdUpdate.Run(); err != nil {
			log.Fatalf("go install failed: %v", err)
		}

	case "com":
		cmdAdd := exec.Command("git", "add", ".")
		stdInsert(cmdAdd)
		if err = cmdAdd.Run(); err != nil {
			log.Fatalf("git add failed: %v", err)
		}	
		
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		cmdCommit := exec.Command("git", "commit", "-m", "don: " + timestamp)
		stdInsert(cmdCommit)
		if err = cmdCommit.Run(); err != nil {
			log.Fatalf("git commit failed: %v", err)
		}
	
		out, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		branch := strings.TrimSpace(string(out))
			
		cmdPush := exec.Command("git", "push", "origin", branch)
		stdInsert(cmdPush)
		if err = cmdPush.Run(); err != nil {
			log.Fatalf("git push failed: %v", err)
		}
	
	case "tail":
		tailPath := os.Args[2]	
		cmdTail := exec.Command("tailwindcss", "-i", tailPath+"/input.css", "-o", tailPath+"/output.css", "--watch")
		stdInsert(cmdTail)
		if err = cmdTail.Run(); err != nil {
			log.Fatalf("tailwindcss failed: %v", err)
		}
		
	default:
		log.Fatalln("Unknown command:", os.Args[1])
	}
}

