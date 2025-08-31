package viewport

type Viewport struct {
	Cursor  int // Current Cursor position in filtered list
	Start   int // Start index for viewport (for scrolling)
	Visible int // Maximum visible items (5 by default)
}
