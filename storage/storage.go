package storage

import "unicorn"

type UnicornStorage interface {
	// Store places a unicorn in storage.
	Store(unicorn *unicorn.Unicorn)

	// InStorage returns the number o unicorns in storage.
	InStorage() int

	// Collect will do a best effort of collecting a number of unicorns from storage.
	// If there are not enough unicorns in storage, it will return any it can provide.
	Collect(n int) []*unicorn.Unicorn
}
