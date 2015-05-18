package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	c := New("test")

	c.Set("foo", "good", -1)

	v, f := c.Get("foo")

	t.Log(v, f)

	c.Set("boy", "sex : nan", 3*time.Second)

	c.Set("boy", "sex : nv", 2*time.Second)

	time.Sleep(2 * time.Second)

	v1, f1 := c.Get("boy")

	t.Logf("%s,%b", v1, f1)

	_ = New("test1")

	t.Logf("%d", len(Client.caches))

	c2 := Client.Get("test")

	c2.Set("foo", "not bad", -1)
	// t.Logf("%p,%p", c, c2)
	v, f = c.Get("foo")
	t.Logf("%s", v)
}
