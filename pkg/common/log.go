package common

import "log/slog"

func LogErrors(errors ...error) {
	for _, err := range errors {
		slog.Warn("Failed indexing", err.Error())
	}
}
