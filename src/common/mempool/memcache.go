package mempool

// A cache holds a set of reusable objects.
// If more are needed, the cache creates them by calling new.

type Cache struct {
	name    string
	newFunc func() interface{}
	bufChan chan interface{}
}

func (c *Cache) Put(x interface{}) {
	select {
	case c.bufChan <- x:
	default:
	}
}
func (c *Cache) Get() (result interface{}) {
	select {
	case result = <-c.bufChan:
	default:
		result = c.newFunc()
	}
	return result
}
func (c *Cache) Count() int {
	return len(c.bufChan)
}
func (c *Cache) Name() string {
	return c.name
}
func (c *Cache) FreeCache() {
	for i := 0; i < len(c.bufChan); i++ {
		<-c.bufChan
	}
}
func NewCache(name string, capacity int, newFunc func() interface{}) *Cache {
	return &Cache{name: name, newFunc: newFunc, bufChan: make(chan interface{}, capacity)}
}
