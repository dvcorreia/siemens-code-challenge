package factory

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"unicorn"
)

//go:embed fixtures/*.txt
var fixtures embed.FS

const (
	defaultName      = "spirit"
	defaultAdjective = "courageous"

	defaultNCapabilities = 3
)

var (
	ErrNotEnoughCapabilities = errors.New("not enough capabilities for producing unicorns")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type factory struct {
	// data for unicorn generation
	names []string
	adj   []string
	cap   []string

	// number of capabilities to attribute to a unicorn
	nCap int
}

// Option is function used to customize the factory.
type Option func(*factory) error

// NCapabilities sets the number of capabilities to attribute to a unicorn.
func NCapabilities(n int) Option {
	return func(f *factory) error {
		if len(f.cap) < n {
			return ErrNotEnoughCapabilities
		}

		f.nCap = n
		return nil
	}
}

// New creates a new unicorn factory.
func New(options ...Option) (*factory, error) {
	names, err := load(fixtures, "fixtures/petnames.txt")
	if err != nil {
		return nil, err
	}

	adj, err := load(fixtures, "fixtures/adj.txt")
	if err != nil {
		return nil, err
	}

	f := &factory{
		names: names,
		adj:   adj,
		cap:   capabilities[:],
		nCap:  defaultNCapabilities,
	}

	for _, opt := range options {
		if opt != nil {
			if err := opt(f); err != nil {
				return nil, err
			}
		}
	}

	if len(f.cap) < f.nCap {
		return nil, ErrNotEnoughCapabilities
	}

	return f, nil
}

// NewUnicorn produces a new unicorn.
func (f factory) NewUnicorn() *unicorn.Unicorn {
	return &unicorn.Unicorn{
		Name:         f.getRandomName(),
		Capabilities: f.selectCapabilities(),
	}
}

// getRandomName generates a random name from the list of adjectives and names.
func (f factory) getRandomName() string {
	name := defaultName
	if len(f.names) != 0 {
		name = f.names[rand.Intn(len(f.names))]
	}

	adj := defaultAdjective
	if len(f.adj) != 0 {
		adj = f.adj[rand.Intn(len(f.adj))]
	}

	return fmt.Sprintf("%s-%s", adj, name)
}

// selectCapabilities select capabilities to give to a unicorn.
func (f factory) selectCapabilities() []string {
	cmap := make(map[string]struct{})

	for i := 0; i < f.nCap; {
		c := f.cap[rand.Intn(len(f.cap))]

		if _, ok := cmap[c]; !ok {
			cmap[c] = struct{}{}
			i++
		}
	}

	caps := make([]string, 0, len(cmap))
	for c := range cmap {
		caps = append(caps, c)
	}

	return caps
}

// load line separated strings from a embeded text file.
func load(fs embed.FS, name string) ([]string, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cap []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		cap = append(cap, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cap, nil
}
