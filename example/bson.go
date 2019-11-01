package example

type BsonInlineType struct{}

type BsonType1 struct {
	BsonInlineType `bson:",inline"`
	Ignore         int `bson:"-"`
	SubKey         int `bson:"-,"`
	//same key
	A              int         `bson:"a,omitempty"`
	InvalidMinSize int         `bson:",minsize"`
	InvalidInline  int         `bson:",inline"`
	InvalidInline2 map[int]int `bson:",inline"`
	ValidInline    struct{}    `bson:",inline"`
}
