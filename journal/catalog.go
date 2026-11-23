package journal

type Catalog[K comparable, V any] struct {
	items map[K]V
	order []V
}

func NewCatalog[K comparable, V any]() *Catalog[K, V] {
	return &Catalog[K, V]{
		items: make(map[K]V),
		order: make([]V, 0),
	}
}

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

func (c *Catalog[K, V]) IsRecorded(key K) bool {
	_, ok := c.items[key]
	return ok
}

func (c *Catalog[K, V]) Existing(key K) (V, bool) {
	item, ok := c.items[key]
	return item, ok
}
