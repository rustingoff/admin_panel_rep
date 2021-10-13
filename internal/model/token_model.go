package model

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}
