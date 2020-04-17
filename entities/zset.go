package entities

type ZSetRequest struct {
	Key   string
	Value interface{}
	Score float64
}

type ZSetResponse struct {
	Value interface{}
	Error error
}
