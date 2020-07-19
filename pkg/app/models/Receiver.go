package models

type ReceiverCache struct {
	ID     uint64 `db:"id"`
	Word   string `db:"word"`
	Answer string `db:"answer"`
	Uses   uint64 `db:"uses"`
}
