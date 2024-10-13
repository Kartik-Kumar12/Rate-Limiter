package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	colorRed = iota + 31
	colorGreen
	colorYellow
	colorMagenta = 35

	colorBold = 1
)

func SetLogger() {

	output := zerolog.ConsoleWriter{Out: os.Stderr}

	// Override the time format
	output.FormatTimestamp = func(i interface{}) string {
		if t, ok := i.(string); ok {
			parsedTime, err := time.Parse(zerolog.TimeFieldFormat, t)
			if err == nil {
				// Custom format for seconds and milliseconds (05.000)
				return parsedTime.Format("05.000")
			}
		}
		// Default if parsing fails
		return time.Now().Format("05.000")
	}

	// Override the colors
	output.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case zerolog.LevelTraceValue:
				l = colorize("TRC", colorMagenta)
			case zerolog.LevelDebugValue:
				l = colorize("DBG", colorMagenta)
			case zerolog.LevelInfoValue:
				l = colorize("INF", colorGreen)
			case zerolog.LevelWarnValue:
				l = colorize(colorize("WRN", colorYellow), colorBold)
			case zerolog.LevelErrorValue:
				l = colorize(colorize("ERR", colorRed), colorBold)
			case zerolog.LevelFatalValue:
				l = colorize(colorize("FTL", colorRed), colorBold)
			case zerolog.LevelPanicValue:
				l = colorize(colorize("PNC", colorRed), colorBold)
			default:
				l = colorize("???", colorBold)
			}
		} else {
			if i == nil {
				l = colorize("???", colorBold)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}

	// Set the global time format for zerolog
	zerolog.TimeFieldFormat = time.RFC3339Nano // Internal storage format for time in logs

	// Assign customized output to the global logger
	log.Logger = log.Output(output)
}

// colorize returns the string s wrapped in ANSI code c
func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
