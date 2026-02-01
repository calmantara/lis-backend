package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type (
	structured struct {
		Level   string `json:"level"`
		Message string `json:"msg"`
	}

	structuredWithPayload struct {
		Level   string `json:"level"`
		Message string `json:"msg"`
		Payload string `json:"payload"`
	}
)

func prepareLoggerMock(name string) (Logger, string) {
	log := logger{}
	config := log.newLoggerConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	fileName := fmt.Sprintf("log_%v.json", uuid.New())
	if name != "" {
		fileName = name
	}
	// logger
	logFile, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	config.Encoding = "json"
	config.ErrorOutputPaths = []string{fileName}
	config.OutputPaths = []string{fileName}
	logger := NewZap(WithZapOption(zap.ErrorOutput(logFile)), WithLoggerConfig(*config))
	return logger, fileName
}

func TestInitZap(t *testing.T) {
	t.Parallel()
	t.Run("init logger", func(t *testing.T) {
		log := NewZap()
		assert.NotNil(t, log)
		// check property
		logs := log.(*logger)
		assert.Equal(t, "lis-backend", logs.applicationName)
		assert.Equal(t, configurations.TEST, logs.environment)
	})
	t.Run("init logger with application name", func(t *testing.T) {
		log := NewZap(WithAppName("lis_backend"))
		assert.NotNil(t, log)
		// check property
		logs := log.(*logger)
		assert.Equal(t, "lis_backend", logs.applicationName)
		assert.Equal(t, configurations.TEST, logs.environment)
	})
	t.Run("init logger with environment", func(t *testing.T) {
		log := NewZap(WithEnvironment(configurations.PRODUCTION), WithInitialFields(map[string]any{
			"type": "worker",
		}))
		assert.NotNil(t, log)
		// check property
		logs := log.(*logger)
		assert.Equal(t, DEFAULT_NAME, logs.applicationName)
		assert.Equal(t, configurations.PRODUCTION, logs.environment)
	})
	t.Run("init logger with config", func(t *testing.T) {
		l := logger{}
		cfg := l.newLoggerConfig()

		log := NewZap(WithLoggerConfig(*cfg))
		assert.NotNil(t, log)
		// check property
		logs := log.(*logger)
		assert.Equal(t, DEFAULT_NAME, logs.applicationName)
		assert.Equal(t, configurations.TEST, logs.environment)
		assert.NotNil(t, logs.cfg)
	})
	t.Run("init logger panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewZap(WithLoggerConfig(zap.Config{}))
		})
	})
}

func TestDebugw(t *testing.T) {
	t.Parallel()
	t.Run("without map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		log.Debugw("debug")

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structured
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "debug")
		assert.Equal(t, result.Message, "debug")
	})

	t.Run("with map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{"payload": "test payload"}
		log.Debugw("debug", additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, "debug", result.Level)
		assert.Equal(t, "debug", result.Message)
		assert.Equal(t, "test payload", result.Payload)
	})
}

func TestDebug(t *testing.T) {
	t.Parallel()

	t.Run("with map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{"payload": "test payload"}
		log.Debug(additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, "debug", result.Level)
		assert.Equal(t, "", result.Message)
		assert.Equal(t, "test payload", result.Payload)
	})
}

func TestInfow(t *testing.T) {
	t.Parallel()
	t.Run("without map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		log.Infow("info")

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structured
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "info")
		assert.Equal(t, result.Message, "info")
	})

	t.Run("with map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{"payload": "test payload"}
		log.Infow("info", additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "info")
		assert.Equal(t, result.Message, "info")
		assert.Equal(t, result.Payload, "test payload")
	})

	t.Run("with empty map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{}
		log.Infow("info", additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "info")
		assert.Equal(t, result.Message, "info")
	})
}

func TestInfo(t *testing.T) {
	t.Parallel()

	t.Run("with map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{"payload": "test payload"}
		log.Info(additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "info")
		assert.Equal(t, result.Message, "")
		assert.Equal(t, result.Payload, "test payload")
	})
}

func TestErrorw(t *testing.T) {
	t.Parallel()
	t.Run("without map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		log.Errorw("error")

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "error")
		assert.Equal(t, result.Message, "error")
	})

	t.Run("with map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{"payload": "test payload"}
		log.Errorw("error", additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "error")
		assert.Equal(t, result.Message, "error")
		assert.Equal(t, result.Payload, "test payload")
	})

	t.Run("with empty map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{}
		log.Errorw("error", additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "error")
		assert.Equal(t, result.Message, "error")
	})
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("with map", func(t *testing.T) {
		log, fileName := prepareLoggerMock("")
		defer os.Remove(fileName)
		additional := map[string]any{"payload": "test payload"}
		log.Error(additional)

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "error")
		assert.Equal(t, result.Message, "")
		assert.Equal(t, result.Payload, "test payload")
	})
}

func TestFatalwWithoutMap(t *testing.T) {
	// Run the crashing code when FLAG is set
	log, fileName := prepareLoggerMock(os.Getenv("TestFatalwWithoutMapFileName"))
	defer os.Remove(fileName)
	if os.Getenv("TestFatalwWithoutMap") == "1" {
		log.Fatalw("fatal")
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalwWithoutMap")
	cmd.Env = append(os.Environ(), "TestFatalwWithoutMap=1", "TestFatalwWithoutMapFileName="+fileName)
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
	// read log file
	dat, _ := os.ReadFile(fileName)
	var result structuredWithPayload
	json.Unmarshal(dat, &result)
	assert.Equal(t, result.Level, "fatal")
	assert.Equal(t, result.Message, "fatal")
}

func TestFatalwWithMap(t *testing.T) {
	// Run the crashing code when FLAG is set
	log, fileName := prepareLoggerMock(os.Getenv("TestFatalwWithMapFileName"))
	defer os.Remove(fileName)
	if os.Getenv("TestFatalwWithMap") == "1" {
		additional := map[string]any{"payload": "test payload"}
		log.Fatalw("fatal", additional)
		return
	}
	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalwWithMap")
	cmd.Env = append(os.Environ(), "TestFatalwWithMap=1", "TestFatalwWithMapFileName="+fileName)
	err := cmd.Run()

	// Cast the error as *exec.ExitError and compare the result
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
	// read log file
	dat, _ := os.ReadFile(fileName)
	var result structuredWithPayload
	json.Unmarshal(dat, &result)
	assert.Equal(t, result.Level, "fatal")
	assert.Equal(t, result.Message, "fatal")
	assert.Equal(t, result.Payload, "test payload")
}

func TestFatal(t *testing.T) {
	t.Parallel()
	t.Run("with map", func(t *testing.T) {
		// Run the crashing code when FLAG is set
		log, fileName := prepareLoggerMock(os.Getenv("TestFatalFileName"))
		defer os.Remove(fileName)
		if os.Getenv("TestFatal") == "1" {
			additional := map[string]any{"payload": "test payload"}
			log.Fatal(additional)
			return
		}
		// Run the test in a subprocess
		cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
		cmd.Env = append(os.Environ(), "TestFatal=1", "TestFatalFileName="+fileName)
		err := cmd.Run()

		// Cast the error as *exec.ExitError and compare the result
		e, ok := err.(*exec.ExitError)
		expectedErrorString := "exit status 1"
		assert.Equal(t, true, ok)
		assert.Equal(t, expectedErrorString, e.Error())

		// read log file
		dat, _ := os.ReadFile(fileName)
		var result structuredWithPayload
		json.Unmarshal(dat, &result)
		assert.Equal(t, result.Level, "fatal")
		assert.Equal(t, result.Message, "")
		assert.Equal(t, result.Payload, "test payload")
	})
}
