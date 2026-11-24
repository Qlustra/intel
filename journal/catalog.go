package journal

// Catalog stores values by key without the visit states of a [Journal]. Its
// zero value is an empty catalog.
type Catalog[K comparable, V any] struct {
	items map[K]V
	order []V
}

// NewCatalog returns an empty catalog.
func NewCatalog[K comparable, V any]() *Catalog[K, V] {
	return &Catalog[K, V]{
		items: make(map[K]V),
		order: make([]V, 0),
	}
}

// Record stores value under key, replacing any previous value for that key.
func (c *Catalog[K, V]) Record(key K, value V) {
	if c.items == nil {
		c.items = make(map[K]V)
	}
	c.items[key] = value
	if c.order == nil {
		c.order = make([]V, 0)
	}
	c.order = append(c.order, value)
}

// IsRecorded reports whether the catalog contains key.
func (c *Catalog[K, V]) IsRecorded(key K) bool {
	_, ok := c.items[key]
	return ok
}

// Existing returns the value for key and whether the key is present.
func (c *Catalog[K, V]) Existing(key K) (V, bool) {
	item, ok := c.items[key]
	return item, ok
}
