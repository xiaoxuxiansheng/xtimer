package xhttp

import "time"

type Option func(*JSONClient)

func WithTimeout(duration time.Duration) Option {
	if duration <= 0 {
		duration = defaultTimeoutDuration
	}
	return func(j *JSONClient) {
		j.timeoutDuration = duration
	}
}

func WithReadLimitBytes(limit int64) Option {
	if limit <= 0 {
		limit = defaultReadLimitBytes
	}

	return func(j *JSONClient) {
		j.readLimitBytes = limit
	}
}

func repair(j *JSONClient) {
	if j.readLimitBytes <= 0 {
		WithReadLimitBytes(defaultReadLimitBytes)(j)
	}

	if j.timeoutDuration <= 0 {
		WithTimeout(defaultTimeoutDuration)(j)
	}
}
