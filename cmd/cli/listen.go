package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forwardURL string

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for webhooks and forward them to a local URL",
	Run: func(cmd *cobra.Command, args []string) {
		gatewayURL := viper.GetString("gateway_url")
		if gatewayURL == "" {
			gatewayURL = "http://localhost:8080"
		}

		u, err := url.Parse(gatewayURL)
		if err != nil {
			log.Fatal(err)
		}

		wsScheme := "ws"
		if u.Scheme == "https" {
			wsScheme = "wss"
		}
		wsURL := fmt.Sprintf("%s://%s/ws", wsScheme, u.Host)

		fmt.Printf("Connecting to %s...\n", wsURL)
		c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		fmt.Printf("Forwarding events to %s\n", forwardURL)
		fmt.Println("Ready! Waiting for events...")

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				fmt.Printf("Received event: %s\n", message)

				// Forward to local URL
				resp, err := http.Post(forwardURL, "application/json", bytes.NewBuffer(message))
				if err != nil {
					fmt.Printf("Failed to forward event: %v\n", err)
				} else {
					fmt.Printf("Forwarded with status: %s\n", resp.Status)
					resp.Body.Close()
				}
			}
		}()

		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("Shutting down...")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			}
		}
	},
}

func init() {
	listenCmd.Flags().StringVarP(&forwardURL, "forward-to", "f", "http://localhost:4242/webhook", "Local URL to forward webhooks to")
	rootCmd.AddCommand(listenCmd)
}
