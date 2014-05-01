package scheme

type Type interface {
	String() string
}

type SchemeType struct {
	Type
	expression string
}
