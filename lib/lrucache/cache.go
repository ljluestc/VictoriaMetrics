package lrucache

// …existing code…

// Get is a backward-compatibility wrapper.
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.GetBytes(key)
}

// Set is a backward-compatibility wrapper.
func (c *Cache) Set(key string, v interface{}) {
	c.SetBytes(key, v)
}
