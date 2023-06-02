package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogging() {

	zerolog.TimeFieldFormat = time.RFC3339

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	logger := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
	logger.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-5s |", i))
	}
	logger.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s,", i)
	}

	log.Logger = zerolog.New(logger).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
