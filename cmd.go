package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()
	paramFlushAll := flag.Bool("all", false, "Use FLUSHALL instead of FLUSHDB")
	paramHost := flag.String("host", "localhost", "Hostname of primary Redis endpoint")
	paramPort := flag.Int("port", 6379, "Port of Redis cluster")
	paramASync := flag.Bool("async", false, "Don't wait for completion")
	paramPassword := flag.String("password", "", "Password for authentication, if enabled")
	paramDB := flag.Int("db", 0, "Index of database to flush (use -all to flush all)")
	flag.Parse()

	if *paramFlushAll && *paramDB != 0 {
		fmt.Fprintf(os.Stderr, "Use -all to target all databases, or -db to specify a specific database other than 0.")
		return
	}

	var res *redis.StatusCmd
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", *paramHost, *paramPort),
		Password: *paramPassword,
		DB:       *paramDB,
	})
	switch {
	case *paramFlushAll && *paramASync:
		res = c.FlushAllAsync(ctx)
	case !*paramFlushAll && *paramASync:
		res = c.FlushDBAsync(ctx)
	case *paramFlushAll && !*paramASync:
		res = c.FlushAll(ctx)
	case !*paramFlushAll && !*paramASync:
		res = c.FlushDB(ctx)
	}
	if res.Err() != nil {
		panic(res.Err())
	}
}
