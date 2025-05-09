package converter

type Converter interface {
	Convert(input, output string) error
}
