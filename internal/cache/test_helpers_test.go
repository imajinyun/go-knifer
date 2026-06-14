package cache

func itoa(i int) string {
	// Simple int-to-string conversion used by tests to avoid extra dependencies.
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	buf := [20]byte{}
	n := 0
	for i > 0 {
		buf[n] = byte('0' + i%10)
		i /= 10
		n++
	}
	if neg {
		buf[n] = '-'
		n++
	}
	// Reverse the digits in place.
	for j, k := 0, n-1; j < k; j, k = j+1, k-1 {
		buf[j], buf[k] = buf[k], buf[j]
	}
	return string(buf[:n])
}

type testTicker struct {
	stopped chan struct{}
}

func (t *testTicker) Stop() { close(t.stopped) }
