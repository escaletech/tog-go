package sessions

import (
	"context"
	"fmt"

	"github.com/twmb/murmur3"

	"github.com/escaletech/tog-go/flags"
)

// NewClient creates a new session client
func NewClient(ctx context.Context, opt ClientOptions) (*Client, error) {
	flagsClient, err := newCachingFlagLister(ctx, flags.ClientOptions{
		Addr:    opt.Addr,
		Cluster: opt.Cluster,
	})
	if err != nil {
		return nil, err
	}

	client := &Client{
		flags:        flagsClient,
		errorHandler: opt.OnError,
	}

	return client, nil
}

type Client struct {
	flags        flagLister
	errorHandler ErrorHandler
}

func (c *Client) Session(ctx context.Context, ns, sid string, opt *SessionOptions) Session {
	opt = initOptions(opt)
	final := opt.Force

	traits := map[string]struct{}{}
	if len(opt.Traits) > 0 {
		for _, t := range opt.Traits {
			traits[t] = struct{}{}
		}
	}

	fs, err := c.flags.ListFlags(ctx, ns)
	if me, ok := err.(flags.MultiError); ok {
		for _, err := range me {
			c.onError(ctx, fmt.Errorf("failed to list flags: %v", err))
		}
	} else if err != nil {
		c.onError(ctx, fmt.Errorf("failed to list flags: %v", err))
		return final
	}

	if len(fs) == 0 {
		return final
	}

	for _, f := range fs {
		if _, isSet := final[f.Name]; isSet {
			continue
		}

		final[f.Name] = resolveFlagValue(f, sid, traits)
	}

	return final
}

func (c *Client) Close() error {
	return c.flags.Close()
}

func (c *Client) onError(ctx context.Context, err error) {
	if c.errorHandler != nil {
		c.errorHandler(ctx, err)
	}
}

func initOptions(opt *SessionOptions) *SessionOptions {
	if opt == nil {
		opt = &SessionOptions{}
	}

	final := &SessionOptions{
		Force:  opt.Force,
		Traits: opt.Traits,
	}
	if final.Force == nil {
		final.Force = Session{}
	}
	if final.Traits == nil {
		final.Traits = []string{}
	}
	return final
}

func resolveFlagValue(f flags.Flag, sid string, traits map[string]struct{}) bool {
	draw := -1
	for _, r := range f.Rollout {
		var pick bool
		if pick, draw = pickRollout(r, sid, traits, f.Timestamp, draw); pick {
			return r.Value
		}
	}

	return false
}

func pickRollout(r flags.Rollout, sid string, traits map[string]struct{}, ts int64, draw int) (bool, int) {
	if r.Percentage != nil {
		if draw < 0 {
			key := []byte(fmt.Sprintf("%v%v", sid, ts))
			hash := murmur3.Sum32(key)
			draw = int(hash % 100)
		}

		if draw >= *r.Percentage {
			return false, draw
		}
	}

	if r.Traits != nil {
		for _, rt := range r.Traits {
			if _, match := traits[rt]; !match {
				return false, draw
			}
		}
	}

	return true, draw
}
