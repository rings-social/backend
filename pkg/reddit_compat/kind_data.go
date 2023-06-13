package reddit_compat

type KindData[T any] struct {
	Kind string `json:"kind"`
	Data T      `json:"data"`
}
