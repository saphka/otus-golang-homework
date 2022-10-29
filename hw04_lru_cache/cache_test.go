package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic simple", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("a", 10)
		require.False(t, wasInCache)
		wasInCache = c.Set("b", 20)
		require.False(t, wasInCache)
		wasInCache = c.Set("c", 30)
		require.False(t, wasInCache)
		wasInCache = c.Set("d", 40)
		require.False(t, wasInCache)

		val, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 20, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 30, val)

		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 40, val)
	})

	t.Run("purge logic read access", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aa", 40)
		require.False(t, wasInCache)
		wasInCache = c.Set("bb", 50)
		require.False(t, wasInCache)
		wasInCache = c.Set("cc", 60)
		require.False(t, wasInCache)

		val, ok := c.Get("bb")
		require.True(t, ok)
		require.Equal(t, 50, val)

		val, ok = c.Get("dd")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("aa")
		require.True(t, ok)
		require.Equal(t, 40, val)

		wasInCache = c.Set("bb", 80)
		require.True(t, wasInCache)

		wasInCache = c.Set("dd", 70)
		require.False(t, wasInCache)

		val, ok = c.Get("cc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
