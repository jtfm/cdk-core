package enums

type Build byte

const (
	Dev Build = iota
	Prod
)

func (build Build) String() string {
	return [...]string{
		"Dev",
		"Prod"}[build]
}
