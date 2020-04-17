package entities

type ZGetRequest struct {
	Key      string
	Value    *string
	ScoreMin float64
	ScoreMax float64
}

type ZGetResponse struct {
	Values []ZGetStruct
	Error  error
}

type ZGetStruct struct {
	Value interface{}
	Score float64
}
