package example

import "io"

type Type3 struct {
	A interface{}   `binding:"exists"`
	B bool          `binding:"exists"`
	C *bool         `binding:"exists"`
	D io.ByteReader `binding:"exists"`
}
