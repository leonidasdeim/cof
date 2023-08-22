package cof

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPutGetPop(t *testing.T) {
	type item struct {
		id  int
		str string
	}

	item1 := item{
		id:  1,
		str: "this is item1",
	}

	item2 := item{
		id:  2,
		str: "and it is item2",
	}

	c, err := Init[item]()
	assert.NoError(t, err)
	defer c.Stop()

	{
		v, ok := c.Get("1")
		require.Empty(t, v)
		require.False(t, ok)
	}

	{
		v, ok := c.Get("2")
		require.Empty(t, v)
		require.False(t, ok)
	}

	c.Put("1", item1)
	c.Put("2", item2)

	{
		v, ok := c.Get("1")
		require.NotEmpty(t, v)
		assert.Equal(t, item1, v)
		assert.True(t, ok)
	}

	{
		v, ok := c.Get("2")
		require.NotEmpty(t, v)
		assert.Equal(t, item2, v)
		assert.True(t, ok)
	}

	{
		v, ok := c.Pop("1")
		require.NotEmpty(t, v)
		assert.Equal(t, item1, v)
		assert.True(t, ok)
	}

	{
		v, ok := c.Pop("2")
		require.NotEmpty(t, v)
		assert.Equal(t, item2, v)
		assert.True(t, ok)
	}

	{
		v, ok := c.Get("1")
		require.Empty(t, v)
		require.False(t, ok)
	}

	{
		v, ok := c.Get("2")
		require.Empty(t, v)
		require.False(t, ok)
	}
}

func TestTTL(t *testing.T) {
	type item struct {
		id  int
		str string
	}

	item1 := item{
		id:  1,
		str: "this is item1",
	}

	c, err := Init[item](TTL(500*time.Millisecond), CleanInterval(50*time.Millisecond))
	assert.NoError(t, err)
	defer c.Stop()

	{
		v, ok := c.Get("1")
		require.Empty(t, v)
		require.False(t, ok)
	}

	c.Put("1", item1)

	{
		v, ok := c.Get("1")
		require.NotEmpty(t, v)
		assert.Equal(t, item1, v)
		assert.True(t, ok)
	}

	time.Sleep(600 * time.Millisecond)

	{
		v, ok := c.Get("1")
		require.Empty(t, v)
		require.False(t, ok)
	}
}
