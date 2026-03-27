package common

import (
	"github.com/charmbracelet/bubbles/key"
)

// ListKeys defines keybindings for the issue list view.
type ListKeys struct {
	Up     key.Binding
	Down   key.Binding
	Filter key.Binding
	Enter  key.Binding
	Close  key.Binding
	Quit   key.Binding
	Help   key.Binding
}

// NewListKeys returns the default list keybindings.
func NewListKeys() ListKeys {
	return ListKeys{
		Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Filter: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "view")),
		Close:  key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "close")),
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

// ShortHelp returns keybindings shown in the compact help view.
func (k ListKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Filter, k.Enter, k.Close, k.Quit}
}

// FullHelp returns keybindings shown in the expanded help view.
func (k ListKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Filter, k.Enter, k.Close},
		{k.Help, k.Quit},
	}
}

// DetailKeys defines keybindings for the issue detail view.
type DetailKeys struct {
	Up   key.Binding
	Down key.Binding
	Back key.Binding
	Quit key.Binding
	Help key.Binding
}

// NewDetailKeys returns the default detail keybindings.
func NewDetailKeys() DetailKeys {
	return DetailKeys{
		Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Back: key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "back")),
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

// ShortHelp returns keybindings shown in the compact help view.
func (k DetailKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back, k.Quit}
}

// FullHelp returns keybindings shown in the expanded help view.
func (k DetailKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Back, k.Quit, k.Help},
	}
}
