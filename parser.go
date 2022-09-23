package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
)

var (
	errHasPrefixEvent = errors.New("hasPrefix event")
	errHasPrefixData  = errors.New("hasPrefix data")

	errCutEvent = errors.New("cut event")
	errCutData  = errors.New("cut data")
)

type event struct {
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

func parse(r io.Reader) ([]event, error) {
	var result []event
	reader := bufio.NewReader(r)
	for {
		var e event

		l, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if err := parseEventName(&e, strings.TrimSpace(string(l))); err != nil {
			log.Println(string(l), err.Error())
			continue
		}

		l, _, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if err := parseEventData(&e, strings.TrimSpace(string(l))); err != nil {
			log.Println(err.Error())
			continue
		}

		result = append(result, e)
	}

	return result, nil
}

func parseEventName(e *event, line string) error {
	var ok bool

	if !strings.HasPrefix(line, "event") {
		return errHasPrefixEvent
	}

	if _, e.Name, ok = strings.Cut(line, ":"); !ok {
		return errCutEvent
	}

	return nil
}

func parseEventData(e *event, line string) error {
	var (
		data string
		ok   bool
	)
	if !strings.HasPrefix(line, "data") {
		return errHasPrefixData
	}

	if _, data, ok = strings.Cut(line, ":"); !ok {
		return errCutData
	}

	if err := json.Unmarshal([]byte(data), &e.Data); err != nil {
		return err
	}

	return nil
}
