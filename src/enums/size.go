package enums

import "database/sql/driver"

type Size string

const (
	XS  Size = "xs"
	SM  Size = "sm"
	MD  Size = "md"
	LG  Size = "lg"
	XL  Size = "xl"
	XXL Size = "xxl"
)

func (s *Size) Scan(value interface{}) error {
	*s = Size(value.(string))
	return nil
}

func (s Size) Value() (driver.Value, error) {
	return string(s), nil
}

func (s Size) String() string {
	return string(s)
}
