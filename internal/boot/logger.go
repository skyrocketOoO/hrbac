package boot

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/serverx/internal/global"
	"github.com/skyrocketOoO/serverx/internal/service/exter/loki"
)

var SendLoki bool

func InitLogger() {
	log.Info().Msg("Logger initialized")

	if global.Env == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if global.Env == "prod" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatTimestamp: func(i interface{}) string {
			if i == nil {
				return "0000-00-00 00:00:00"
			}
			return i.(string)
		},
		FormatLevel: func(i interface{}) string {
			if i == nil {
				return "[???]"
			}
			return "[" + i.(string) + "]"
		},
		FormatCaller: func(i interface{}) string {
			if i == nil {
				return "unknown:0"
			}
			return simplifyCaller(i.(string))
		},
		FormatMessage: func(i interface{}) string {
			if i == nil {
				return ""
			}
			return i.(string)
		},
		// NoColor: false,
	}

	log.Info().Msg("Logger initialized")
	if SendLoki {
		consoleWriter.Out = log.Output(loki.NewLokiWriter()).With().Caller().Timestamp().Logger()
		log.Logger = log.Output(loki.NewLokiWriter()).With().Caller().Timestamp().Logger()
	} else {
		log.Logger = log.Output(consoleWriter).With().Caller().Timestamp().Logger()
	}
}

func simplifyCaller(caller string) string {
	file := filepath.Base(caller)
	dir := filepath.Dir(caller)

	return filepath.Join(filepath.Base(dir), file)
}
