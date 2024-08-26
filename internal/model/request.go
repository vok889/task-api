package model

import "task-api/internal/constant"

type RequestItem struct {
	Title    string  ``
	Price    float64 `binding:"required,gt=0"`
	Quantity uint
}

type RequestFindItem struct {
	Statuses []constant.ItemStatus `form:"status[]"`
}

type RequestUpdateItem struct {
	Status constant.ItemStatus
}
