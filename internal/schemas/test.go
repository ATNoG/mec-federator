package schemas

type Test struct {
	Number int    `bson:"number,omitempty"`
	Text   string `bson:"text,omitempty"`
}
