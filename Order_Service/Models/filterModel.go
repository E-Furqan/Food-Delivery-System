package model

type Filter struct {
	ColumnName     string `json:"column_name"`
	OrderDirection string `json:"order_direction"`
}
