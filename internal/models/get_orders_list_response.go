package models

type GetOrdersListResponse []OrderModel

func NewGetOrdersListResponse(orders []OrderModel) *GetOrdersListResponse {
	return (*GetOrdersListResponse)(&orders)
}
