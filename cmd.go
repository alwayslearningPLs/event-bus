package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/r3labs/sse"
	"github.com/spf13/cobra"
)

var (
	url        string
	inputFile  string
	outputFile string
)

func newRootCmd() *cobra.Command {
	var root = &cobra.Command{
		Version: "0.0.1",
		Use:     "event_bus",
		Short:   "event bus filtering",
		Long:    "event bus filtering",
		PreRun: func(cmd *cobra.Command, args []string) {
			if strings.TrimSpace(outputFile) == "" {
				if strings.TrimSpace(inputFile) != "" {
					outputFile = inputFile + ".json"
				} else {
					outputFile = "events.json"
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if inputFile != "" {
				dumpFile()
			} else {
				var events []event
				cl := sse.NewClient(url)

				if b, err := os.ReadFile(outputFile); err == nil {
					if err = json.Unmarshal(b, &events); err != nil {
						os.Remove(outputFile)
					}
				}

				log.Println("Listening events on ", url)
				cl.SubscribeRaw(func(msg *sse.Event) {
					e := event{Name: string(msg.Event)}
					if err := json.Unmarshal(msg.Data, &e.Data); err != nil {
						log.Println("error trying to unmarshal the event: ", msg.Data)
						return
					}

					events = append(events, e)

					b, err := json.Marshal(e)
					if err != nil {
						log.Println("error trying to marshal the event: ", e)
						return
					}

					os.WriteFile(outputFile, b, 0666)
				})
			}
		},
	}

	return root
}

func execute() error {
	var root = newRootCmd()

	root.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "event streams that we are going to parse. We can only specify an inputFile and the output file will be the same name plus '.json'")
	root.PersistentFlags().StringVarP(&url, "url", "u", "", "url to listen from server event streams")
	root.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "event streams that we are going to marshal")

	root.MarkPersistentFlagFilename("input-file")
	root.MarkPersistentFlagFilename("output-file")

	return root.Execute()
}

func dumpFile() {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	events, err := parse(f)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(events, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(outputFile, b, 0666); err != nil {
		log.Fatal(err)
	}
}
