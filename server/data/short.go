package data

type ShortInfo struct {
	Id       string `json:"id,onitempty" bson:"id,onitempty"`             // The id of the short url
	Redirect string `json:"redirect,onitempty" bson:"redirect,onitempty"` // The url to redirect to
	Uses     int    `json:"uses,onitempty" bson:"uses,onitempty"`         // The amount of uses
}
