package generator

import "github.com/teris-io/shortid"

type URLGenerator struct{}

func (URLGenerator) Generate() (string, error) {
	return shortid.Generate()
}
