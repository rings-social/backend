package response

type Paginated[T any] struct {
	Items []T `json:"items"`

	// The total number of items.
	Total int64 `json:"total"`

	After string `json:"after"`
}
