package maryhelp

import (
	"github.com/sirupsen/logrus"
	"github.com/ansel1/merry"
)

// Hook struct
type Hook struct {
	stacktraceType StacktraceType
}

// StacktraceType set how display stacktrace
type StacktraceType int

const (
	// NONE not show stacktrace
	NONE StacktraceType = iota
	// TEXT set stacktrace format as text, raw from merry
	TEXT
	// JSON set stacktrace format as json (map)
	JSON
)

// NewHook create a new hook
func NewHook(stacktraceType StacktraceType) Hook {
	return Hook{
		stacktraceType: stacktraceType,
	}
}

// Levels method
func (hook Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func buildJSONMerryLog(err merry.Error) map[string]interface{} {

	stacktrace, flag := merryStacktraceJSON(err)
	if flag != 0 {
		return nil
	}

	return map[string]interface{}{
		"stacktrace": stacktrace,
	}
}

func buildTextMerryLog(err merry.Error) map[string]interface{} {

	return map[string]interface{}{
		"stacktrace": merry.Stacktrace(err),
	}
}

// Fire method
func (hook Hook) Fire(entry *logrus.Entry) error {

	if entry.Data["error"] == nil {
		return nil
	}

	switch entry.Data["error"].(type) {
	case merry.Error:
		err := entry.Data["error"].(merry.Error)
		delete(entry.Data, "error")

		switch hook.stacktraceType {
		case NONE:
			return nil
		case TEXT:
			merryData := buildTextMerryLog(err)
			merryAppendStringValues(merryData, err)
			entry.Data["merry"] = merryData
			return nil
		case JSON:
			merryData := buildJSONMerryLog(err)
			merryAppendStringValues(merryData, err)
			entry.Data["merry"] = merryData
			return nil
		}
	}

	return nil
}
