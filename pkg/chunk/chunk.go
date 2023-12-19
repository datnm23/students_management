package chunk

func Chunk[T any](s []T, n int) [][]T {
	var chunk [][]T
	for i := 0; i < len(s); i += n {
		step := i + n

		if step > len(s) {
			step = len(s)
		}
		chunk = append(chunk, s[i:step])
	}
	return chunk
}
