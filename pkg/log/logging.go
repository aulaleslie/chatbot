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

	file, err := os.OpenFile(
		"chatbot.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	logger := zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339, NoColor: true}
	logger.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-5s |", i))
	}
	logger.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s,", i)
	}

	multiWriter := zerolog.MultiLevelWriter(logger, file)
	writer := zerolog.SyncWriter(multiWriter)

	log.Logger = zerolog.New(writer).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
