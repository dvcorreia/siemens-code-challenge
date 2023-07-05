package storage

import (
	"log"
	"unicorn"
)

type storageLogger struct {
	logger *log.Logger
	store  UnicornStorage
}

// WithLogs adds logging to a unicorn storage.
func WithLogs(logger *log.Logger, store UnicornStorage) UnicornStorage {
	return &storageLogger{
		logger: logger,
		store:  store,
	}
}

// Store places a unicorn in storage.
func (l *storageLogger) Store(unicorn *unicorn.Unicorn) {
	l.store.Store(unicorn)
	l.logger.Printf("storage: stored unicorn<%s>, now with %d", unicorn.Name, l.store.InStorage())
}

// InStorage returns the number o unicorns in storage.
func (l *storageLogger) InStorage() int {
	return l.store.InStorage()
}

// Collect will do a best effort of collecting a number of unicorns from storage.
// If there are not enough unicorns in storage, it will return any it can provide.
func (l *storageLogger) Collect(n int) []*unicorn.Unicorn {
	unicorns := l.store.Collect(n)
	l.logger.Printf("storage: collected %d from the requested %d", len(unicorns), n)
	return unicorns
}
