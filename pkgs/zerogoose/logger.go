package zerogoose

import "github.com/rs/zerolog"

type gooseZeroLogger struct {
	log zerolog.Logger
}

func New(logger zerolog.Logger) gooseZeroLogger {
	return gooseZeroLogger{
		log: logger,
	}
}

func (l gooseZeroLogger) Fatal(v ...any) {
	l.log.Fatal().Msgf("%v", v...)
}

func (l gooseZeroLogger) Fatalf(format string, v ...any) {
	l.log.Fatal().Msgf(format, v...)
}

func (l gooseZeroLogger) Print(v ...any) {
	l.log.Info().Msgf("%v", v...)
}

func (l gooseZeroLogger) Printf(format string, v ...any) {
	l.log.Info().Msgf(format, v...)
}
