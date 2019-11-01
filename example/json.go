package example

type errorCode int

type Type1 struct {
	A int `json:"a"`
	B int `json:"B"`
	C int `json:"C,omitempty"`
	D int `json:",string"`
	E struct {
		A int
	} `json:"e,string"`
	F           int        `json:"-"`
	G           *int       `json:"-,string"`
	H           *errorCode `json:",string"`
	I           errorCode  `json:",string"`
	InvalidName errorCode  `json:"\\,string"`
}

type Type2 struct {
	D struct {
		E int `json:"E"`
		F int `json:"F,omitempty"`
	}
}
