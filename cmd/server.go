// Copyright © 2017 Alexander Pinnecke <alexander.pinnecke@googlemail.com>

package cmd

import (
	"log"
	"net/http"

	"github.com/StageAutoControl/visualizer/command"
	"github.com/StageAutoControl/visualizer/websocket"
	"github.com/spf13/cobra"
)

var (
	frontendListen string
	commandListen  string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the visualizer server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server := websocket.NewServer()
		go server.Run()

		receiver := command.NewReceiver(commandListen)
		receiver.AddHandler(server)
		go receiver.Receive()
		log.Printf("Command server listening on %s \n", commandListen)

		mux := http.NewServeMux()
		mux.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
			server.ServeRequest(w, r)
		})

		log.Printf("Websocket server listening on %s \n", frontendListen)
		err := http.ListenAndServe(frontendListen, mux)
		if err != nil {
			log.Fatalf("Failed to listen on frontend http: %v \n", err)
		}

	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&frontendListen, "frontend-port", "f", "0.0.0.0:3001", "The listen string to bind the frontend server to")
	serverCmd.Flags().StringVarP(&commandListen, "command-port", "c", "0.0.0.0:1337", "The listen string to bind the command receiver to")

}

func errorHandler() chan error {
	c := make(chan error, 1)

	go func() {
		for {
			err := <-c
			log.Println(err)
		}
	}()

	return c
}
