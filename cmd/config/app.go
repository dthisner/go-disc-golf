package config

import "log/slog"

// Not implemented - Closures for dependency injection
// Look at: https://gist.github.com/alexedwards/5cd712192b4831058b21
type Application struct {
	Logger *slog.Logger
}
