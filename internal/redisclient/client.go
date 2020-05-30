package redisclient

import "github.com/go-redis/redis/v8"

func New(addr string, cluster bool) (redis.UniversalClient, error) {
	opt, err := redis.ParseURL(addr)
	if err != nil {
		return nil, err
	}

	opt.MaxRetries = 3

	if !cluster {
		return redis.NewClient(opt), nil
	}

	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{opt.Addr},
	}), nil
}
