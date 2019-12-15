package test_test

import (
	"testing"

	"github.com/lucmski/limo/extras/httpcache"
	"github.com/lucmski/limo/extras/httpcache/test"
)

func TestMemoryCache(t *testing.T) {
	test.Cache(t, httpcache.NewMemoryCache())
}
