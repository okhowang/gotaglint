package example

type Type1 struct {
	A int `json:"a"`
	B int `json:"B"`
	C int `json:"C,omitempty"`
}

type Type2 struct {
	D struct {
		E int `json:"E"`
		F int `json:"F,omitempty"`
	}
}
