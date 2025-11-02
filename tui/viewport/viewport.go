package viewport

type Viewport struct {
	Cursor  int // Current Cursor position in filtered list
	Offset  int // Offset index for viewport (for scrolling)
	Visible int // Maximum visible items (5 by default)
}
