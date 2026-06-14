package vref

type facadeSample struct {
	Name   string `json:"name"`
	hidden string
}

func (s facadeSample) GetName() string      { return s.Name }
func (s facadeSample) Add(a int, b int) int { return a + b }

type facadeEmbedded struct {
	Code string `ref:"code"`
}

type facadeExtendedSample struct {
	facadeEmbedded
	Alias string `xml:"alias"`
	Count int
}

type facadeMethodSample struct{}

func (facadeMethodSample) Equal(facadeMethodSample) bool { return true }
func (facadeMethodSample) HashCode() int                 { return 7 }
func (facadeMethodSample) String() string                { return "method-sample" }
func (facadeMethodSample) SetName(string)                {}

func newFacadeSample(name string) facadeSample { return facadeSample{Name: name} }

type facadeError struct{}

func (facadeError) Error() string { return "facade" }
