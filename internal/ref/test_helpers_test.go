package ref

type embeddedSample struct {
	Base string `json:"base"`
}

type sample struct {
	embeddedSample
	Name   string `json:"name"`
	Age    int
	hidden string
}

func (s sample) GetName() string         { return s.Name }
func (s *sample) SetName(name string)    { s.Name = name }
func (s sample) Add(a int, b int) int    { return a + b }
func (s sample) String() string          { return s.Name }
func (s sample) Equal(other sample) bool { return s.Name == other.Name }
func (s sample) HashCode() int           { return len(s.Name) }

func newSample(name string, age int) sample { return sample{Name: name, Age: age} }

type sampleError struct{}

func (sampleError) Error() string { return "sample" }

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
