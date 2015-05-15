package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	c := New()

	c.Set("foo", "good", 0)

	v, f := c.Get("foo")

	t.Log(v, f)

	c.Set("boy", "sex : nan", 3*time.Second)

	c.Set("boy", "sex : nv", 2*time.Second)

	time.Sleep(2 * time.Second)

	v1, f1 := c.Get("boy")

	t.Logf("%s,%b", v1, f1)
}
