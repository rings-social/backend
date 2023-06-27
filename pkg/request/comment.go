package request

type Comment struct {
	Content  string `json:"content"`
	ParentId *uint  `json:"parentId"`
}
