package reddit_compat

type Listing[T any] struct {
	After     *string       `json:"after"`
	Before    *string       `json:"before"`
	Dist      int           `json:"dist"`
	ModHash   string        `json:"modhash"`
	GeoFilter string        `json:"geo_filter"`
	Children  []KindData[T] `json:"children,nilasempty"`
}
