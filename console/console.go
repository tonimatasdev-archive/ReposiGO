package console

import (
	"bufio"
	"github.com/TonimatasDEV/ReposiGO/session"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Console(server *http.Server) {
	go func() {
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-stopChan:
			handleStop(server)
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		rawCommand, err := inputReader.ReadString('\n')

		if err != nil {
			log.Println("Exception on read the command:", err)
			continue
		}

		command := strings.Replace(rawCommand, "\n", "", -1)
		split := strings.Split(command, " ")

		label := split[0]

		switch label {
		case "quit", "exit", "stop":
			handleStop(server)
		case "help", "?":
			handleHelp()
		case "session":
			handleSession(split)
		default:
			if command != "" {
				log.Println("Unknown command \"" + label + "\".")
			}
		}
	}
}

func handleStop(server *http.Server) {
	log.Println("ReposiGO stopped successfully.")
	_ = server.Close()
	os.Exit(0)
}

func handleHelp() {
	log.Print("ReposiGO commands: \n" +
		"	\"stop\",\"exit\" or \"quit\" to stop the server.\n" +
		"	\"session create <username> <read-access> <write-access>\" to create a session.\n" +
		"	\"session delete <username>\" to delete a session.\n")
}

func handleSession(split []string) {
	if len(split) < 2 {
		log.Println("Subcommand is needed. Available subcommands are \"create\" and \"delete\".")
		return
	}

	switch split[1] {
	case "create":
		if len(split) < 5 {
			log.Println("More args are needed: session create <username> <read-access> <write-access>.\n" +
				"	Read/write access can be \"NONE\", \"REPO1,REPO2,REPO3\" or \"*\" for all.")
			return
		}

		readAccess := strings.Split(split[3], ",")
		writeAccess := strings.Split(split[4], ",")

		session.CreateSession(split[2], readAccess, writeAccess)
	case "delete":
		if len(split) < 3 {
			log.Println("More args are needed: session delete <username>.")
			return
		}

		session.DeleteSession(split[2])
	default:
		log.Println("Unknown session subcommand: " + split[1] + ".\n" +
			"	Available subcommands are \"create\" and \"delete\".")
	}
}
