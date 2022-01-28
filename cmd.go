package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var paramFlushAll bool
var paramHost string
var paramPort int
var paramASync bool
var paramPassword string
var paramDB int

func init() {
	flag.BoolVar(&paramFlushAll, "all", false, "Use FLUSHALL instead of FLUSHDB")
	flag.StringVar(&paramHost, "host", "localhost", "Hostname of primary Redis endpoint")
	flag.IntVar(&paramPort, "port", 6379, "Port of Redis cluster")
	flag.BoolVar(&paramASync, "async", false, "Don't wait for completion")
	flag.StringVar(&paramPassword, "password", "", "Password for authentication, if enabled")
	flag.IntVar(&paramPort, "db", 0, "Index of database to flush (use -all to flush all)")
	flag.Parse()
}

func main() {
	var res *redis.StatusCmd
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", paramHost, paramPort),
		Password: paramPassword,
		DB:       paramDB,
	})
	switch {
	case paramFlushAll && paramASync:
		res = c.FlushAllAsync(context.Background())
	case !paramFlushAll && paramASync:
		res = c.FlushDBAsync(context.Background())
	case paramFlushAll && !paramASync:
		res = c.FlushAll(context.Background())
	case !paramFlushAll && !paramASync:
		res = c.FlushDB(context.Background())
	}
	if res.Err() != nil {
		panic(res.Err())
	}
}
