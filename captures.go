package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"time"
)

func captureOption(title string, options []string, prompt string, allowEmpty bool, reader *bufio.Reader) int {
	color.Blue(title)
	for i := 1; i <= len(options); i++ {
		color.Cyan("%v) %v", i, options[i-1])
	}
	fmt.Printf("\n%v: ", prompt)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.Replace(text, "\n", "", -1))

	if text == "" {
		if allowEmpty {
			return -1
		} else {
			color.Red("ERR: invalid option\n\n")
			return captureOption(title, options, prompt, allowEmpty, reader)
		}
	} else {
		option, err := strconv.Atoi(text)
		if err != nil {
			color.Red("ERR: invalid option\n\n")
			return captureOption(title, options, prompt, allowEmpty, reader)
		}
		if option > 0 && option <= len(options) {
			return option - 1
		} else {
			color.Red("ERR: invalid option\n\n")
			return captureOption(title, options, prompt, allowEmpty, reader)
		}
	}
}

func captureText(prompt string, defaultAnswer string, reader *bufio.Reader) string {
	c := color.New(color.FgCyan)
	c.DisableColor()
	_, _ = c.Print(prompt)
	if defaultAnswer != "" {
		c.EnableColor()
		_, _ = c.Printf(" (%v)", defaultAnswer)
		c.DisableColor()
	}
	_, _ = c.Print(": ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.Replace(text, "\n", "", -1))
	if text == "" {
		if defaultAnswer == "" {
			color.Red("ERR: invalid input\n")
			return captureText(prompt, defaultAnswer, reader)
		}
		return defaultAnswer
	}
	return text
}

func captureDate(prompt string, defaultDate time.Time, reader *bufio.Reader) time.Time {
	c := color.New(color.FgCyan)
	c.DisableColor()
	_, _ = c.Print(prompt)
	c.EnableColor()
	_, _ = c.Printf(" (%v/%v)", int(defaultDate.Month()), defaultDate.Day())
	c.DisableColor()
	_, _ = c.Print(": ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.Replace(text, "\n", "", -1))
	if text == "" {
		return defaultDate
	}
	timeInput, err := time.Parse("01/02", text)
	if err != nil {
		color.Red("ERR: invalid format, must be MM/dd")
		return captureDate(prompt, defaultDate, reader)
	}
	return timeInput
}

func captureBool(prompt string, reader *bufio.Reader) bool {
	c := color.New(color.FgCyan)
	c.DisableColor()
	_, _ = c.Print(prompt)
	c.EnableColor()
	_, _ = c.Print(" (y for yes, n for no)")
	c.DisableColor()
	_, _ = c.Print("?: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.Replace(text, "\n", "", -1))

	if text == "y" {
		return true
	} else if text == "n" {
		return false
	} else {
		color.Red("ERR: invalid input")
		return captureBool(prompt, reader)
	}
}

func captureMoney(prompt string, defaultAnswer int, reader *bufio.Reader) int {
	c := color.New(color.FgCyan)
	c.DisableColor()
	_, _ = c.Print(prompt)
	if defaultAnswer >= 0 {
		c.EnableColor()
		_, _ = c.Printf(" ($%v)", float64(defaultAnswer) / 100.0)
		c.DisableColor()
	}
	_, _ = c.Print(": ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.Replace(text, "\n", "", -1))
	if text == "" {
		if defaultAnswer < 0 {
			color.Red("ERR: invalid input")
			return captureMoney(prompt, defaultAnswer, reader)
		}
		return defaultAnswer
	}
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		color.Red("ERR: invalid input")
		return captureMoney(prompt, defaultAnswer, reader)
	}

	return int(f * 100)
}