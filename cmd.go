package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

type logger bool

func (l logger) stdout(f string, a ...interface{}) {
	if !l {
		return
	}
	fmt.Fprintf(os.Stderr, f+"\n", a...)
}

func main() {
	ctx := context.Background()
	paramDryRun := flag.Bool("dryrun", false, "Only attempt a no-op PING command")
	paramVerbose := flag.Bool("verbose", false, "Output more text")
	paramFlushAll := flag.Bool("all", false, "Use FLUSHALL instead of FLUSHDB")
	paramHost := flag.String("host", "localhost", "Hostname of primary Redis endpoint")
	paramPort := flag.Int("port", 6379, "Port of Redis cluster")
	paramASync := flag.Bool("async", false, "Don't wait for completion")
	paramPassword := flag.String("password", "", "Password for authentication, if enabled")
	paramDB := flag.Int("db", 0, "Index of database to flush (use -all to flush all)")
	flag.Parse()

	log := logger(*paramVerbose)
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
	log.stdout("Will connect to: %s", c.Options().Addr)
	if *paramPassword != "" {
		log.stdout("Using password for connection.")
	}
	switch {
	case *paramDryRun:
		log.stdout("Sending PING")
		res = c.Ping(ctx)
	case *paramFlushAll && *paramASync:
		log.stdout("Sending FLUSHALL ASYNC")
		res = c.FlushAllAsync(ctx)
	case !*paramFlushAll && *paramASync:
		log.stdout("Sending FLUSHDB ASYNC (database %d)", *paramDB)
		res = c.FlushDBAsync(ctx)
	case *paramFlushAll && !*paramASync:
		log.stdout("Sending FLUSHALL")
		res = c.FlushAll(ctx)
	case !*paramFlushAll && !*paramASync:
		log.stdout("Sending FLUSHDB (database %d)", *paramDB)
		res = c.FlushDB(ctx)
	}
	if res.Err() != nil {
		panic(res.Err())
	}
	log.stdout("Command result: %+v", res)
}
