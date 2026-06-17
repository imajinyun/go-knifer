package obj

type sample struct {
	Name string
	Tags []string
}

type timeLike struct {
	Value string
}

type registeredValue struct {
	Value string
}

type namedRegisteredValue struct {
	Value string
}

func mustPanic(t testingT, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

type testingT interface {
	Helper()
	Fatal(args ...any)
}
