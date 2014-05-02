package scheme

type Type interface {
	String() string
	IsNumber() bool
}

type SchemeType struct {
	expression string
}

func (s *SchemeType) IsNumber() bool {
	return false
}
