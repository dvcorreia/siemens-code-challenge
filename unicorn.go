package unicorn

// RequestID is used to identify pending unicorn production requests.
type RequestID string

// Unicorn is a horse with a beautiful horn.
// They are have funny names and can do a lot of stuff.
type Unicorn struct {
	Name         string
	Capabilities []string
}

// UnicornService is a service that can produce happy beautiful unicorns.
type UnicornService interface {
	// RequestUnicorns initiates a new unicorn production request.
	// If no sufficient unicorn are available, it returns a request ID for consequent pooling.
	RequestUnicorns() (RequestID, []Unicorn)

	// Pool returns avalilable unicorns for a request ID order.
	Pool(RequestID) []Unicorn

	// Validate check is a request ID has an orden in the process.
	Validate(RequestID) bool
}
