package maryhelp

import (
	"strings"

	"github.com/ansel1/merry"
)

func isOdd(n int) bool {
	return n%2 != 0
}

func merryAppendStringValues(data map[string]interface{}, err merry.Error) {

	for key, value := range merry.Values(err) {
		switch key.(type) {
		case string:
			data[key.(string)] = value
		}
	}
}

func merryStacktraceJSON(err merry.Error) ([]map[string]string, int) {

	r := []map[string]string{}

	stacktrace := merry.Stacktrace(err)
	l := strings.Split(stacktrace, "\n")

	n := len(l)
	if isOdd(n) {
		n--
	}

	if n < 2 {
		return nil, 1
	}

	for i := 0; i < n; i += 2 {

		// parse line 0
		fileLineAddress := strings.Split(l[i], " ")
		if len(fileLineAddress) < 2 {
			return nil, 1
		}
		fileLine := strings.Split(fileLineAddress[0], ":")
		if len(fileLine) < 2 {
			return nil, 2
		}

		// parse line 1
		funcCode := strings.Split(l[i+1], ":")
		if len(funcCode) < 2 {
			return nil, 3
		}

		r = append(r, map[string]string{
			"file":    strings.Trim(fileLine[0], " "),
			"line":    strings.Trim(fileLine[1], " "),
			"address": strings.Trim(fileLineAddress[1], " ()"),
			"func":    strings.Trim(funcCode[0], " \t"),
			"code":    strings.Trim(funcCode[1], " "),
		})
	}

	return r, 0
}
