package cli

import (
	"fmt"
	"time"
)

type Level int

const (
	FATAL Level = iota + 1
	ERROR
	WARNING
	DEBUG
	INFO
	TRACE
)

type Cli struct {
	Level Level
}

func (c *Cli) PrintTrace(msg string) { c.print(TRACE, "TRACE", msg) }

func (c *Cli) PrintInfo(msg string) { c.print(INFO, "INFO", msg) }

func (c *Cli) PrintDebug(msg string) { c.print(DEBUG, "DEBUG", msg) }

func (c *Cli) PrintWarning(msg string) { c.print(WARNING, "WARNING", msg) }

func (c *Cli) PrintError(msg string) { c.print(ERROR, "ERROR", msg) }

func (c *Cli) PrintFatal(msg string) { c.print(FATAL, "FATAL", msg) }

func (c *Cli) print(level Level, prefix string, msg string) {
	if c.Level < level {
		return
	}

	timeNow := time.Now()
	fmt.Printf("[%s][%s]%s\n",
		timeNow.Format("2006-01-02 15:04:05"),
		prefix,
		msg,
	)
}
