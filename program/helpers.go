package program

func IsWindowSizeValid(width, height int) bool {
	return width >= minWindowWidth && height >= minWindowHeight
}
