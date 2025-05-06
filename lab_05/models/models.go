package models

type Stop struct {
	StopId   int    `json:"stopId"`
	StopName string `json:"stopName"`
	StopCode string `json:"stopCode"`
}

type Departure struct {
	RouteID       int    `json:"routeId"`
	HeadSign      string `json:"headSign"`
	EstimatedTime string `json:"estimatedTime"`
}
