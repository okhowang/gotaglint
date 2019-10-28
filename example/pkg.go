package example

import "io"

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

type Type3 struct {
	A interface{}   `binding:"exists"`
	B bool          `binding:"exists"`
	C *bool         `binding:"exists"`
	D io.ByteReader `binding:"exists"`
}
