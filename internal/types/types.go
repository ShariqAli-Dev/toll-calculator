package types

const (
	KAFKA_TOPIC         = "obudata"
	AGGREGATOR_ENDPOINT = "http://127.0.0.1:3001"
)

// should be transport independant
type Invoice struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
}

type Distance struct {
	Values float64 `json:"value"`
	OBUID  int     `json:"obuID"`
	Unix   int64   `json:"unix"`
}

type OBUData struct {
	OBUID     int     `json:"obuID"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	RequestID int     `json:"requestID"`
}
