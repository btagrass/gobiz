package mdl

type Cache struct {
	Id          int64  `json:"id"`
	Key         string `json:"key"`
	Value       any    `json:"value"`
	ValueString string `json:"valueString"`
	Expiration  int64  `json:"expiration"`
}
