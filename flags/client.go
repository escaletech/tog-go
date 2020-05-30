package flags

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/escaletech/tog-go/internal/keys"
	"github.com/escaletech/tog-go/internal/redisclient"
)

func NewClient(opt ClientOptions) (*Client, error) {
	r, err := redisclient.New(opt.Addr, opt.Cluster)
	if err != nil {
		return nil, err
	}

	return &Client{
		redis: r,
	}, nil
}

type Client struct {
	redis redis.UniversalClient
}

func (c *Client) ListFlags(ctx context.Context, ns string) ([]Flag, error) {
	res, err := c.redis.HGetAll(ctx, keys.Flags(ns)).Result()
	if err != nil {
		return nil, err
	}

	flags := make([]Flag, len(res))
	i := 0
	errs := MultiError{}
	for name, raw := range res {
		f, err := c.parseFlag(ns, name, raw)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		flags[i] = f
		i++
	}

	if len(errs) > 0 {
		return nil, errs
	}

	sort.Slice(flags, func(i, j int) bool {
		return strings.Compare(flags[i].Name, flags[j].Name) < 0
	})

	return flags, nil
}

func (c *Client) GetFlag(ctx context.Context, ns, name string) (Flag, error) {
	raw, err := c.redis.HGet(ctx, keys.Flags(ns), name).Result()
	if err != nil {
		return Flag{}, err
	}

	return c.parseFlag(ns, name, raw)
}

func (c *Client) SaveFlag(ctx context.Context, f Flag) (Flag, error) {
	if err := validate(f); err != nil {
		return Flag{}, err
	}

	sanitized := f
	sanitized.Name = ""
	sanitized.Namespace = ""
	sanitized.Timestamp = time.Now().Unix()
	f.Timestamp = sanitized.Timestamp

	raw, err := json.Marshal(sanitized)
	if err != nil {
		return Flag{}, err
	}

	if err := c.redis.HSet(ctx, keys.Flags(f.Namespace), f.Name, string(raw)).Err(); err != nil {
		return Flag{}, err
	}

	if err := c.redis.Publish(ctx, keys.PubSub, f.Namespace).Err(); err != nil {
		return Flag{}, err
	}

	return f, nil
}

func (c *Client) DeleteFlag(ctx context.Context, ns, name string) error {
	if err := c.redis.HDel(ctx, keys.Flags(ns), name).Err(); err != nil {
		return err
	}

	if err := c.redis.Publish(ctx, keys.PubSub, ns).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Close() error {
	if err := c.redis.Close(); err != nil {
		return fmt.Errorf("failed to close redis client: %v", err)
	}

	return nil
}

func (c *Client) parseFlag(ns, name, raw string) (Flag, error) {
	var f Flag
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		return Flag{}, err
	}

	f.Namespace = ns
	f.Name = name

	return f, nil
}
