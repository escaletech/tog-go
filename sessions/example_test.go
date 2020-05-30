package sessions_test

import (
	"context"
	"fmt"
	"math"

	"github.com/go-redis/redis/v8"

	"github.com/escaletech/tog-go/flags"
	"github.com/escaletech/tog-go/sessions"
)

var ctx = context.Background()
var fc *flags.Client
var sc *sessions.Client
var rdb *redis.Client

func init() {
	opt := sessions.ClientOptions{
		Addr:    "redis://localhost:6379/2",
		Cluster: false,
	}

	redisOpt, _ := redis.ParseURL(opt.Addr)
	rdb = redis.NewClient(redisOpt)
	rdb.FlushDB(ctx)

	var err error
	sc, err = sessions.NewClient(ctx, opt)
	if err != nil {
		panic(err)
	}

	fc, err = flags.NewClient(flags.ClientOptions{
		Addr:    opt.Addr,
		Cluster: opt.Cluster,
	})
	if err != nil {
		panic(err)
	}
}

func ExampleNewClient() {
	sc, err := sessions.NewClient(ctx, sessions.ClientOptions{
		Addr:    "redis://localhost:6379",
		Cluster: false,
	})
	defer sc.Close()

	fmt.Println(err)
	// Output: <nil>
}

func ExampleClient_Session_one_flag() {
	fc.SaveFlag(ctx, flags.Flag{
		Namespace: "my_app_simple",
		Name:      "blue-button",
		Rollout:   []flags.Rollout{{Value: true}},
	})

	sess := sc.Session(ctx, "my_app_simple", "session_id", nil)

	fmt.Println(sess)
	// Output: map[blue-button:true]
}

func ExampleClient_Session_empty() {
	sess := sc.Session(ctx, "my_app_empty", "session_id", nil)

	fmt.Println(sess)
	// Output: map[]
}

func ExampleClient_Session_random() {
	const reps = 100000
	percentage := 50
	fc.SaveFlag(ctx, flags.Flag{
		Namespace: "my_app_random",
		Name:      "blue-button",
		Rollout:   []flags.Rollout{{Value: true, Percentage: &percentage}},
	})

	var trueCount, falseCount float64
	for i := 0; i < reps; i++ {
		sess := sc.Session(ctx, "my_app_random", fmt.Sprintf("session_id_%v", i), nil)
		if sess["blue-button"] {
			trueCount++
		} else {
			falseCount++
		}
	}

	trueRate := math.Round(trueCount / (reps / 100))
	falseRate := math.Round(falseCount / (reps / 100))
	fmt.Println(trueRate, "/", falseRate)
	// Output: 50 / 50
}

func ExampleClient_Session_override_value() {
	fc.SaveFlag(ctx, flags.Flag{
		Namespace: "my_app_override",
		Name:      "blue-button",
		Rollout:   []flags.Rollout{{Value: false}},
	})

	sess := sc.Session(ctx, "my_app_override", "session_id", &sessions.SessionOptions{
		Force: sessions.Session{"blue-button": true},
	})

	fmt.Println(sess)
	// Output: map[blue-button:true]
}

func ExampleClient_Session_value_by_traits() {
	fc.SaveFlag(ctx, flags.Flag{
		Namespace: "my_app_traits",
		Name:      "blue-button",
		Rollout: []flags.Rollout{
			{Value: true, Traits: []string{"early-adopter"}},
			{Value: false},
		},
	})

	earlyAdopterSession := sc.Session(ctx, "my_app_traits", "session_id", &sessions.SessionOptions{
		Traits: []string{"early-adopter"},
	})

	commonSession := sc.Session(ctx, "my_app_traits", "session_id", nil)

	fmt.Println(earlyAdopterSession)
	fmt.Println(commonSession)
	// Output: map[blue-button:true]
	// map[blue-button:false]
}
