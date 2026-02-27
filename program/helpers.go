package program

func IsWindowSizeValid(width, height int) bool {
	return width >= minWindowWidth && height >= minWindowHeight
}

func ptr[T any](v T) *T { return &v }
