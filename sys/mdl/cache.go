package mdl

type Cache struct {
	Id         int64  `json:"id"`
	Key        string `json:"key"`
	Val        any    `json:"val"`
	ValString  string `json:"valString"`
	Expiration int64  `json:"expiration"`
}
