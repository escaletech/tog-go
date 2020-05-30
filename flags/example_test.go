package flags_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/escaletech/tog-go/flags"
	"github.com/escaletech/tog-go/internal/keys"
)

var ctx = context.Background()
var fc *flags.Client
var rdb *redis.Client

func init() {
	opt := flags.ClientOptions{
		Addr:    "redis://localhost:6379/2",
		Cluster: false,
	}

	redisOpt, _ := redis.ParseURL(opt.Addr)
	rdb = redis.NewClient(redisOpt)
	rdb.FlushDB(ctx)

	var err error
	fc, err = flags.NewClient(opt)
	if err != nil {
		panic(err)
	}

	rdb.HMSet(ctx, keys.Flags("ns_with_some_flags"), map[string]interface{}{
		"blue-button": `{"description":"blue_descr", "timestamp": 1, "rollout": [{"value": true}]}`,
		"white-bg":    `{"description":"white_descr", "timestamp": 2, "rollout": [{"percentage": 30, "value": true}]}`,
	})

	rdb.HSet(ctx, keys.Flags("ns_with_one_flag"),
		"blue-button",
		`{"description":"blue_descr", "timestamp": 1, "rollout": [{"value": true}]}`)
}

func ExampleNewClient() {
	fc, err := flags.NewClient(flags.ClientOptions{
		Addr:    "redis://localhost:6379/2",
		Cluster: false,
	})
	defer fc.Close()

	fmt.Println(err)
	// Output: <nil>
}

func ExampleClient_ListFlags() {
	fs, err := fc.ListFlags(ctx, "ns_with_some_flags")

	asJSON, _ := json.MarshalIndent(fs, "", "  ")
	fmt.Println("error:", err)
	fmt.Println(string(asJSON))

	// Output: error: <nil>
	// [
	//   {
	//     "namespace": "ns_with_some_flags",
	//     "name": "blue-button",
	//     "description": "blue_descr",
	//     "timestamp": 1,
	//     "rollout": [
	//       {
	//         "value": true
	//       }
	//     ]
	//   },
	//   {
	//     "namespace": "ns_with_some_flags",
	//     "name": "white-bg",
	//     "description": "white_descr",
	//     "timestamp": 2,
	//     "rollout": [
	//       {
	//         "value": true,
	//         "percentage": 30
	//       }
	//     ]
	//   }
	// ]
}

func ExampleClient_ListFlags_Empty() {
	fs, err := fc.ListFlags(ctx, "ns_empty")

	asJSON, _ := json.MarshalIndent(fs, "", "  ")
	fmt.Println("error:", err)
	fmt.Println(string(asJSON))

	// Output: error: <nil>
	// []
}

func ExampleClient_GetFlag() {
	fc.SaveFlag(ctx, flags.Flag{Namespace: "ns_get_flag", Name: "my-flag"})

	f, err := fc.GetFlag(ctx, "ns_get_flag", "my-flag")
	fmt.Println(err)
	fmt.Println(f.Name)

	// Output: <nil>
	// my-flag
}

func ExampleClient_SaveFlag_and_GetFlag() {
	f, err := fc.SaveFlag(ctx, flags.Flag{
		Namespace:   "my_app_save_flag",
		Name:        "blue-button",
		Description: "some_description",
		Rollout:     []flags.Rollout{{Value: true}},
	})

	fmt.Println("returned:", err, f.Name)
	f2, _ := fc.GetFlag(ctx, "my_app_save_flag", "blue-button")
	fmt.Println("found:", f2.Name)

	// Output: returned: <nil> blue-button
	// found: blue-button
}

func ExampleClient_DeleteFlag() {
	before, _ := fc.ListFlags(ctx, "ns_with_one_flag")
	fmt.Println("before count:", len(before))

	err := fc.DeleteFlag(ctx, "ns_with_one_flag", "blue-button")

	fmt.Println("error:", err)

	after, _ := fc.ListFlags(ctx, "ns_with_one_flag")
	fmt.Println("after count:", len(after))

	// Output: before count: 1
	// error: <nil>
	// after count: 0
}
