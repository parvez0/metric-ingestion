package objects

import "time"

type GenericResponse struct {
	Success bool `json="success"`
	Message string `json="message"`
	Data interface{} `json:"data"`
}

type Metrics struct {
	RowId int64 `json:"row_id"`
	CpuUsed int `json:"percentage_cpu_used"`
	MemoryUsed int `json:"percentage_memory_used"'`
	Ip string `json:"ip"`
	Date time.Time `json:"date"`
}

type MetricsList []Metrics

type ApiError struct {
	Error string
	Status int
}