package keys

const PubSub = "tog3:namespace-changed"

func Flags(ns string) string {
	return "tog3:flags:" + ns
}
