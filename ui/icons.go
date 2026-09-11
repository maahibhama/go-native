package ui

import (
	"strings"
	"unicode"
)

// Standard cross-platform icon names.
const (
	IconChevronRight = "chevron.right"
	IconChevronLeft  = "chevron.left"
	IconChevronUp    = "chevron.up"
	IconChevronDown  = "chevron.down"
	IconArrowLeft    = "arrow.left"
	IconArrowRight   = "arrow.right"
	IconArrowUp      = "arrow.up"
	IconArrowDown    = "arrow.down"
	IconCheck        = "check"
	IconClose        = "close"
	IconPlus         = "plus"
	IconMinus        = "minus"
	IconSearch       = "search"
	IconSettings     = "settings"
	IconStar         = "star"
	IconHeart        = "heart"
	IconShare        = "share"
	IconTrash        = "trash"
	IconRefresh      = "refresh"
	IconInfo         = "info"
	IconWarning      = "warning"
	IconError        = "error"
	IconHome         = "home"
	IconLock         = "lock"
	IconUnlock       = "unlock"
	IconEdit         = "edit"
	IconFilter       = "filter"
)

var standardIcons = []string{
	IconChevronRight,
	IconChevronLeft,
	IconChevronUp,
	IconChevronDown,
	IconArrowLeft,
	IconArrowRight,
	IconArrowUp,
	IconArrowDown,
	IconCheck,
	IconClose,
	IconPlus,
	IconMinus,
	IconSearch,
	IconSettings,
	IconStar,
	IconHeart,
	IconShare,
	IconTrash,
	IconRefresh,
	IconInfo,
	IconWarning,
	IconError,
	IconHome,
	IconLock,
	IconUnlock,
	IconEdit,
	IconFilter,
}

var standardIconSet = map[string]struct{}{
	IconChevronRight: {},
	IconChevronLeft:  {},
	IconChevronUp:    {},
	IconChevronDown:  {},
	IconArrowLeft:    {},
	IconArrowRight:   {},
	IconArrowUp:      {},
	IconArrowDown:    {},
	IconCheck:        {},
	IconClose:        {},
	IconPlus:         {},
	IconMinus:        {},
	IconSearch:       {},
	IconSettings:     {},
	IconStar:         {},
	IconHeart:        {},
	IconShare:        {},
	IconTrash:        {},
	IconRefresh:      {},
	IconInfo:         {},
	IconWarning:      {},
	IconError:        {},
	IconHome:         {},
	IconLock:         {},
	IconUnlock:       {},
	IconEdit:         {},
	IconFilter:       {},
}

// StandardIcons returns a slice of all recognized standard cross-platform icon names.
func StandardIcons() []string {
	return append([]string(nil), standardIcons...)
}

// IsStandardIcon reports whether name is one of the recognized standard cross-platform icon names.
func IsStandardIcon(name string) bool {
	_, ok := standardIconSet[name]
	return ok
}

// IconAccessibleLabel generates a human-readable accessibility label for an icon name.
func IconAccessibleLabel(name string) string {
	switch name {
	case IconChevronRight:
		return "Chevron right"
	case IconChevronLeft:
		return "Chevron left"
	case IconChevronUp:
		return "Chevron up"
	case IconChevronDown:
		return "Chevron down"
	case IconArrowLeft:
		return "Arrow left"
	case IconArrowRight:
		return "Arrow right"
	case IconArrowUp:
		return "Arrow up"
	case IconArrowDown:
		return "Arrow down"
	case IconCheck:
		return "Check"
	case IconClose:
		return "Close"
	case IconPlus:
		return "Plus"
	case IconMinus:
		return "Minus"
	case IconSearch:
		return "Search"
	case IconSettings:
		return "Settings"
	case IconStar:
		return "Star"
	case IconHeart:
		return "Heart"
	case IconShare:
		return "Share"
	case IconTrash:
		return "Trash"
	case IconRefresh:
		return "Refresh"
	case IconInfo:
		return "Info"
	case IconWarning:
		return "Warning"
	case IconError:
		return "Error"
	case IconHome:
		return "Home"
	case IconLock:
		return "Lock"
	case IconUnlock:
		return "Unlock"
	case IconEdit:
		return "Edit"
	case IconFilter:
		return "Filter"
	default:
		return formatFallbackIconLabel(name)
	}
}

func formatFallbackIconLabel(name string) string {
	if name == "" {
		return ""
	}
	cleaned := strings.Map(func(r rune) rune {
		if r == '.' || r == '_' || r == '-' {
			return ' '
		}
		return r
	}, name)
	parts := strings.Fields(cleaned)
	if len(parts) == 0 {
		return ""
	}
	runes := []rune(parts[0])
	runes[0] = unicode.ToUpper(runes[0])
	parts[0] = string(runes)
	return strings.Join(parts, " ")
}

// PlatformIcon resolves a standard cross-platform icon name to its native symbol or drawable name.
func PlatformIcon(name string, platform Platform) string {
	if platform == PlatformIOS {
		switch name {
		case IconChevronRight:
			return "chevron.right"
		case IconChevronLeft:
			return "chevron.left"
		case IconChevronUp:
			return "chevron.up"
		case IconChevronDown:
			return "chevron.down"
		case IconArrowLeft:
			return "arrow.left"
		case IconArrowRight:
			return "arrow.right"
		case IconArrowUp:
			return "arrow.up"
		case IconArrowDown:
			return "arrow.down"
		case IconCheck:
			return "checkmark"
		case IconClose:
			return "xmark"
		case IconPlus:
			return "plus"
		case IconMinus:
			return "minus"
		case IconSearch:
			return "magnifyingglass"
		case IconSettings:
			return "gearshape"
		case IconStar:
			return "star"
		case IconHeart:
			return "heart"
		case IconShare:
			return "square.and.arrow.up"
		case IconTrash:
			return "trash"
		case IconRefresh:
			return "arrow.clockwise"
		case IconInfo:
			return "info.circle"
		case IconWarning:
			return "exclamationmark.triangle"
		case IconError:
			return "xmark.circle"
		case IconHome:
			return "house"
		case IconLock:
			return "lock"
		case IconUnlock:
			return "lock.open"
		case IconEdit:
			return "pencil"
		case IconFilter:
			return "line.3.horizontal.decrease.circle"
		default:
			return name
		}
	}
	if platform == PlatformAndroid {
		switch name {
		case IconChevronRight:
			return "ic_menu_forward"
		case IconChevronLeft:
			return "ic_menu_revert"
		case IconCheck:
			return "checkbox_on_background"
		case IconClose:
			return "ic_menu_close_clear_cancel"
		case IconPlus:
			return "ic_input_add"
		case IconMinus:
			return "ic_input_delete"
		case IconSearch:
			return "ic_menu_search"
		case IconSettings:
			return "ic_menu_preferences"
		case IconStar:
			return "star_on"
		case IconShare:
			return "ic_menu_share"
		case IconTrash:
			return "ic_menu_delete"
		case IconRefresh:
			return "ic_popup_sync"
		case IconInfo:
			return "ic_dialog_info"
		case IconWarning:
			return "ic_dialog_alert"
		case IconError:
			return "stat_notify_error"
		default:
			return name
		}
	}
	return name
}

// Icon creates a cross-platform icon component using NodeImage, configured with accessible
// role RoleImage and a default human-readable accessible label.
func Icon(name string) *element {
	e := newElement(NodeImage, Props{
		ImageSource: name,
		ImageMode:   ImageFit,
		AccessRole:  RoleImage,
		AccessLabel: IconAccessibleLabel(name),
	})
	return e.IconSize(24)
}

// IconSize sets the icon's width and height in points.
func (e *element) IconSize(points float32) *element {
	if points < 0 {
		points = 0
	}
	return e.Width(points).Height(points)
}

// IconColor sets the icon's foreground tint color.
func (e *element) IconColor(color Color) *element {
	e.node.Style.Appearance.Foreground = color
	e.node.Style.Text.Color = color
	return e
}
