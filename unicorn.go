package unicorn

// Unicorn is a horse with a beautiful horn.
// They are have funny names and can do a lot of stuff.
type Unicorn struct {
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
}

// OrderID is used to identify pending unicorn production request orders.
type OrderID string

// Service is a service that can produce happy beautiful unicorns.
type Service interface {
	// RequestUnicorns initiates a new unicorn production request.
	// If no sufficient unicorn are available, it returns a request ID for consequent pooling.
	OrderUnicorns(n int) (OrderID, error)

	// Pool returns the available ordered unicorns and how many are left to produce.
	Pool(OrderID) (int, []Unicorn)

	// Validate checks if an ID has an orden in the process.
	Validate(OrderID) bool
}
