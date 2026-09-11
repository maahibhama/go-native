package ui

import "strconv"

// Checkbox creates an accessible checkbox control.
func Checkbox(checked bool, onChange func(bool)) *element {
	e := newElement(NodeSwitch, Props{Checked: checked, AccessRole: RoleButton})
	e.node.Toggle = onChange
	e.node.Press = func() {
		if !e.node.Style.Interaction.Disabled && onChange != nil {
			onChange(!checked)
		}
	}
	applyButtonTouchTargets(e)
	return e
}

// Radio creates an accessible single-choice radio element.
func Radio(selected bool, onSelect func()) *element {
	e := newElement(NodeButton, Props{Checked: selected, AccessRole: RoleButton})
	e.node.Toggle = func(b bool) {
		if b && !e.node.Style.Interaction.Disabled && onSelect != nil {
			onSelect()
		}
	}
	e.node.Press = func() {
		if !selected && !e.node.Style.Interaction.Disabled && onSelect != nil {
			onSelect()
		}
	}
	applyButtonTouchTargets(e)
	return e
}

// RadioOption specifies a choice in a RadioGroup.
type RadioOption[T comparable] struct {
	Value    T
	Label    string
	Disabled bool
}

type radioGroupComponent[T comparable] struct {
	selected T
	onSelect func(T)
	options  []RadioOption[T]
}

// RadioGroup presents mutually exclusive choices in a column.
func RadioGroup[T comparable](selected T, onSelect func(T), options []RadioOption[T]) Component {
	return &radioGroupComponent[T]{
		selected: selected,
		onSelect: onSelect,
		options:  append([]RadioOption[T](nil), options...),
	}
}

func (r *radioGroupComponent[T]) Build() *Node {
	return r.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (r *radioGroupComponent[T]) BuildContext(context BuildContext) *Node {
	rows := make([]Component, len(r.options))
	for i, opt := range r.options {
		val := opt.Value
		isSelected := val == r.selected
		selectFn := func() {
			if !opt.Disabled && r.onSelect != nil {
				r.onSelect(val)
			}
		}
		radio := Radio(isSelected, selectFn)
		if opt.Disabled {
			radio.Disabled(true)
		}
		var item Component
		if opt.Label != "" {
			label := Text(opt.Label)
			if opt.Disabled {
				label.Disabled(true)
			}
			row := Row(radio, label).Gap(8).Align(AlignCenter)
			if opt.Disabled {
				row.Disabled(true)
			}
			item = row
		} else {
			item = radio
		}
		rows[i] = item
	}
	col := Column(rows...).Gap(8)
	return col.BuildContext(context)
}

// Slider creates a continuous or stepped value adjustment control.
func Slider(value float32, onChange func(float32)) *element {
	e := newElement(NodeProgressIndicator, Props{Progress: value, AccessRole: RoleButton})
	e.sliderMin = 0.0
	e.sliderMax = 1.0
	e.sliderStep = 0.0
	e.sliderChange = onChange
	e.node.SliderMin = 0.0
	e.node.SliderMax = 1.0
	e.node.SliderStep = 0.0
	e.node.SliderChange = onChange
	e.node.Change = func(v string) {
		if f, err := strconv.ParseFloat(v, 32); err == nil && onChange != nil {
			onChange(float32(f))
		}
	}
	applyButtonTouchTargets(e)
	return e
}

// SliderMin sets the minimum value of a Slider.
func (e *element) SliderMin(min float32) *element {
	e.sliderMin = min
	e.node.SliderMin = min
	return e
}

// SliderMax sets the maximum value of a Slider.
func (e *element) SliderMax(max float32) *element {
	e.sliderMax = max
	e.node.SliderMax = max
	return e
}

// SliderStep sets the discretization step of a Slider.
func (e *element) SliderStep(step float32) *element {
	e.sliderStep = step
	e.node.SliderStep = step
	return e
}

// SegmentItem specifies one option in a SegmentedControl.
type SegmentItem[T comparable] struct {
	Value    T
	Label    string
	Icon     string
	Disabled bool
}

type segmentedControlComponent[T comparable] struct {
	selected T
	onSelect func(T)
	items    []SegmentItem[T]
}

// SegmentedControl presents a horizontal single-selection control.
func SegmentedControl[T comparable](selected T, onSelect func(T), items []SegmentItem[T]) Component {
	return &segmentedControlComponent[T]{
		selected: selected,
		onSelect: onSelect,
		items:    append([]SegmentItem[T](nil), items...),
	}
}

func (s *segmentedControlComponent[T]) Build() *Node {
	return s.BuildContext(NewBuildContext(DefaultEnvironment()))
}

func (s *segmentedControlComponent[T]) BuildContext(context BuildContext) *Node {
	theme := context.Environment.Theme
	primaryColor := theme.Colors["primary"]
	if primaryColor == (Color{}) {
		primaryColor = RGB(0, 122, 255)
	}
	segments := make([]Component, len(s.items))
	for i, item := range s.items {
		val := item.Value
		isSelected := val == s.selected
		selectFn := func() {
			if !item.Disabled && s.onSelect != nil {
				s.onSelect(val)
			}
		}
		label := item.Label
		if label == "" && item.Icon != "" {
			label = item.Icon
		}
		btn := Button(label, selectFn)
		if isSelected {
			btn.Background(primaryColor).Foreground(RGB(255, 255, 255)).Bold()
		} else {
			btn.Background(RGBA(0, 0, 0, 0)).Foreground(primaryColor)
		}
		btn.PaddingXY(12, 6)
		btn.CornerRadius(6)
		if item.Disabled {
			btn.Disabled(true)
		}
		segments[i] = btn
	}
	container := Row(segments...).
		Border(1, primaryColor).
		CornerRadius(8).
		Padding(2).
		Align(AlignCenter)
	return container.BuildContext(context)
}

// LinearProgress creates a linear progress indicator. Negative value indicates indeterminate progress.
func LinearProgress(value float32) *element {
	progress := value
	if progress >= 0 {
		if progress > 1 {
			progress = 1
		}
	} else {
		progress = -1
	}
	e := newElement(NodeProgressIndicator, Props{Progress: progress})
	e.node.IsCircular = false
	return e
}

// CircularProgress creates a circular progress indicator. Negative value indicates indeterminate progress.
func CircularProgress(value float32) *element {
	progress := value
	if progress >= 0 {
		if progress > 1 {
			progress = 1
		}
	} else {
		progress = -1
	}
	e := newElement(NodeProgressIndicator, Props{Progress: progress})
	e.node.IsCircular = true
	return e
}
