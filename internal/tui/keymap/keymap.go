package keymap

type Action int

const (
	ActionNone Action = iota
	ActionQuit
	ActionHelp
	ActionSearch
	ActionBlur
	ActionCopy
	ActionDown
	ActionUp
	ActionLeft
	ActionRight
	ActionTop
	ActionBottom
	ActionHalfDown
	ActionHalfUp
	ActionPageDown
	ActionPageUp
	ActionHistoryNext
	ActionHistoryPrev
)

type Resolver struct {
	pendingG bool
}

func (r *Resolver) Resolve(key string, editing bool) Action {
	if key == "" {
		return ActionNone
	}

	if r.pendingG && key != "g" {
		r.pendingG = false
	}
	if editing {
		switch key {
		case "?":
			return ActionHelp
		case "esc":
			return ActionBlur
		case "ctrl+c":
			return ActionQuit
		case "ctrl+p", "N":
			return ActionHistoryPrev
		case "ctrl+n", "n":
			return ActionHistoryNext
		case "enter":
			return ActionCopy
		default:
			return ActionNone
		}
	}

	switch key {
	case "ctrl+c", "q":
		return ActionQuit
	case "?":
		return ActionHelp
	case "/":
		return ActionSearch
	case "esc":
		return ActionBlur
	case "enter":
		return ActionCopy
	case "j", "down":
		return ActionDown
	case "k", "up":
		return ActionUp
	case "h", "left":
		return ActionLeft
	case "l", "right":
		return ActionRight
	case "G":
		return ActionBottom
	case "ctrl+d":
		return ActionHalfDown
	case "ctrl+u":
		return ActionHalfUp
	case "ctrl+f", "pgdown":
		return ActionPageDown
	case "ctrl+b", "pgup":
		return ActionPageUp
	case "n":
		return ActionHistoryNext
	case "N":
		return ActionHistoryPrev
	case "g":
		if r.pendingG {
			r.pendingG = false
			return ActionTop
		}
		r.pendingG = true
		return ActionNone
	default:
		return ActionNone
	}
}
