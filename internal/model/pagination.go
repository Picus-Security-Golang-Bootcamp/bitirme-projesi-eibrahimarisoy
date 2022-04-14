package model

type Search struct {
	Column string `json:"column"`
	Action string `json:"action"`
	Query  string `json:"query"`
}

// pagination model
// type Pagination struct {
// 	Limit        int         `json:"limit"`
// 	Page         int         `json:"page"`
// 	Sort         string      `json:"sort"`
// 	TotalRows    int64       `json:"total_rows"`
// 	FirstPage    string      `json:"first_page"`
// 	PreviousPage string      `json:"previous_page"`
// 	NextPage     string      `json:"next_page"`
// 	LastPage     string      `json:"last_page"`
// 	Rows         interface{} `json:"rows"`
// 	Q            string      `json:"q"`
// }
