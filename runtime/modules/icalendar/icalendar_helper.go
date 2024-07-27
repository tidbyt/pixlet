package icalendar

import (
	"bufio"
	"net/http"
)

type ICalendar struct {
	url  string
	data *bufio.Scanner
}

func NewICalendar(url string) *ICalendar {
	return &ICalendar{
		url:  url,
		data: nil,
	}
}

func (c *ICalendar) GetCalendar() error {
	data, err := http.Get(c.url)
	if err != nil {
		return err
	}

	c.data = bufio.NewScanner(data.Body)

	return err
}

func (c *ICalendar) ParseCalendar() error {
	return nil
}
