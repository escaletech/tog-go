package sessions

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	lru "github.com/hashicorp/golang-lru"

	"github.com/escaletech/tog-go/flags"
	"github.com/escaletech/tog-go/internal/keys"
	"github.com/escaletech/tog-go/internal/redisclient"
)

func newCachingFlagLister(ctx context.Context, opt flags.ClientOptions) (*cachingFlagLister, error) {
	flagsClient, err := flags.NewClient(opt)
	if err != nil {
		return nil, err
	}

	subscriber, err := redisclient.New(opt.Addr, opt.Cluster)
	if err != nil {
		return nil, err
	}

	cache, err := lru.New(10)
	if err != nil {
		return nil, err
	}

	lister := &cachingFlagLister{
		flags:      flagsClient,
		subscriber: subscriber,
		cache:      cache,
	}

	if err := lister.init(ctx); err != nil {
		return nil, err
	}

	return lister, nil
}

type cachingFlagLister struct {
	flags      flagLister
	subscriber redis.UniversalClient
	cache      *lru.Cache
}

func (c *cachingFlagLister) init(ctx context.Context) error {
	pubsub := c.subscriber.Subscribe(ctx, keys.PubSub)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}

	ch := pubsub.Channel()
	go func() {
		for m := range ch {
			c.cache.Remove(m.Payload)
		}
	}()

	return nil
}

func (c *cachingFlagLister) ListFlags(ctx context.Context, ns string) ([]flags.Flag, error) {
	value, isCached := c.cache.Get(ns)
	if isCached {
		fs, ok := value.([]flags.Flag)
		if !ok {
			return nil, fmt.Errorf("unexpected value %v in cache key", value)
		}

		return fs, nil
	}

	fs, err := c.flags.ListFlags(ctx, ns)
	if err != nil {
		return nil, fmt.Errorf("failed to list flags: %v", err)
	}

	c.cache.Add(ns, fs)
	return fs, nil
}

func (c *cachingFlagLister) Close() error {
	if err := c.flags.Close(); err != nil {
		return fmt.Errorf("failed to close flag lister: %v", err)
	}

	if err := c.subscriber.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber: %v", err)
	}

	return nil
}
