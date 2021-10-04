package comfig

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	sentryhook "github.com/xr9kayu/logrus/sentry"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

type Logger interface {
	Log() *logan.Entry
}

type logger struct {
	getter  kv.Getter
	once    Once
	options LoggerOpts
}

type LoggerOpts struct {
	Release string
}

func NewLogger(getter kv.Getter, options LoggerOpts) Logger {
	return &logger{
		getter:  getter,
		options: options,
	}
}

func (l *logger) Log() *logan.Entry {
	return l.once.Do(func() interface{} {
		config, err := parseLogConfig()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out logger"))
		}

		entry := logan.New().Level(config.Level)

		if !config.DisableSentry {
			sentrier := NewSentrier(l.getter, SentryOpts{Release: l.options.Release})
			sentry := sentrier.Sentry()
			sentryConfig := sentrier.SentryConfig()

			var selectedLevel logrus.Level
			if sentryConfig.Level != nil {
				selectedLevel = logrus.Level(*sentryConfig.Level)
			} else {
				selectedLevel = logrus.Level(config.Level)
			}
			levels := make([]logrus.Level, 0)
			for level := logrus.PanicLevel; level <= selectedLevel; level++ {
				levels = append(levels, level)
			}
			hook, err := sentryhook.NewHook(sentry, levels...)
			if err != nil {
				panic(errors.Wrap(err, "failed to init sentry hook"))
			}

			entry.AddLogrusHook(hook)
		}

		return entry
	}).(*logan.Entry)
}

type loggerConfig struct {
	Level         logan.Level 
	DisableSentry bool
}

func parseLogConfig() (*loggerConfig, error) {
	config := loggerConfig{
		Level: logan.ErrorLevel,
	}

	var err error 
	envResult, isSet := os.LookupEnv("LOG_LEVEL")
	if(isSet) {
		config.Level, err = logan.ParseLevel(envResult)
		if err != nil {
			panic("LOG_LEVEL not set")
		}
	} else {
		panic("LOG_LEVEL no set")
	}
	envResult, isSet = os.LookupEnv("LOG_DISABLE_SENTRY")
	if(isSet) {
		config.DisableSentry, err = strconv.ParseBool(envResult)
		if err != nil {
			panic("LOG_LEVEL not set")
		}
	} else {
		panic("LOG_DISABLE_SENTRY no set")
	}

	return &config, nil
}
