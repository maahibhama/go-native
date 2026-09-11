package ui

// Spacer expands along its parent's main axis. Use FixedSpacer when a
// deterministic gap is required instead of flexible remaining space.
func Spacer() *element { return View().Flex(1, 1) }

// FixedSpacer creates equal logical width and height. In a Row only width is
// normally relevant; in a Column only height is normally relevant.
func FixedSpacer(size float32) *element { return View().Width(size).Height(size) }

// Center composes flex spacers around child so it is centered on both axes
// whenever its parent gives it bounded space.
func Center(child Component) Component {
	return Column(
		Spacer(),
		Row(Spacer(), child, Spacer()).Flex(1, 1).Align(AlignCenter),
		Spacer(),
	).Flex(1, 1).Align(AlignCenter)
}

// AspectRatio constrains child through a wrapper with the requested width to
// height ratio. A width or height constraint on the wrapper resolves the other.
func AspectRatio(ratio float32, child Component) *element {
	return View(child).AspectRatio(ratio)
}

// Stack overlays children in declaration order. The first child participates
// in intrinsic sizing; subsequent children are absolutely positioned. Give
// the Stack an explicit or parent-derived size when every child is an overlay.
func Stack(children ...Component) Component {
	return &stackComponent{children: append([]Component(nil), children...)}
}

type stackComponent struct{ children []Component }

func (s *stackComponent) Build() *Node {
	return s.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (s *stackComponent) BuildContext(context BuildContext) *Node {
	root := View().BuildContext(context)
	for index, child := range s.children {
		if child == nil {
			continue
		}
		node := BuildWithContext(child, context.Child("stack:"+itoa(index)))
		if node == nil {
			continue
		}
		if index > 0 {
			node.Style.Layout.Position = PositionAbsolute
		}
		root.Children = append(root.Children, node)
	}
	return root
}

// KeyboardAvoidingView adds the current bottom keyboard inset to its padding.
// Native hosts update MediaQuery when the keyboard geometry changes.
func KeyboardAvoidingView(child Component) Component {
	return &keyboardAvoidingComponent{child: child}
}

type keyboardAvoidingComponent struct{ child Component }

func (k *keyboardAvoidingComponent) Build() *Node {
	return k.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (k *keyboardAvoidingComponent) BuildContext(context BuildContext) *Node {
	insets := context.Environment.MediaQuery.KeyboardInsets
	return View(k.child).Styled(Style{Layout: LayoutStyle{Padding: EdgeInsets{Bottom: insets.Bottom}}}).BuildContext(context)
}

// Divider creates a one-point horizontal separator using the supplied color.
func Divider(color Color) *element { return View().Height(1).Background(color) }

// VerticalDivider creates a one-point vertical separator.
func VerticalDivider(color Color) *element { return View().Width(1).Background(color) }

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var buffer [20]byte
	index := len(buffer)
	for value > 0 {
		index--
		buffer[index] = byte('0' + value%10)
		value /= 10
	}
	return string(buffer[index:])
}
