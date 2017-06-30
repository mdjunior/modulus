package handlers

import cache "github.com/patrickmn/go-cache"

// DBHandler is a struct that
type DBHandler struct {
	cache *cache.Cache
}

// NewDBHandler return a new pointer of db struct
func NewDBHandler(cache *cache.Cache) *DBHandler {
	return &DBHandler{
		cache: cache,
	}
}
