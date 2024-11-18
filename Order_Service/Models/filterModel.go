package model

type Filter struct {
	ColumnName string `json:"column_name"`
	SortOrder  string `json:"sort_order"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
