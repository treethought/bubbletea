package tea

import (
	"context"
	"image/color"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Context represents a Bubble Tea program's context. It is passed to the
// program's Init, Update, and View functions to provide information about the
// program's state and to allow them to interact with the terminal.
type Context interface {
	context.Context

	// SetValue sets a value on the context. This is useful for storing values
	// that needs to be accessed across multiple functions.
	// You can access the value later using Value.
	SetValue(key, value interface{})

	// BackgroundColor returns the current background color of the terminal.
	// It returns nil if the terminal's doesn't support querying the background
	// color.
	BackgroundColor() color.Color

	// HasLightBackground returns true if the terminal's background color is
	// light. This is useful for determining whether to use light or dark colors
	// in the program's UI.
	HasLightBackground() bool

	// SupportsEnhancedKeyboard reports whether the terminal supports enhanced
	// keyboard keys. On Windows, this means it supports virtual keys like and
	// the Windows Console API. On Unix, this means it supports the Kitty
	// Keyboard Protocol.
	SupportsEnhancedKeyboard() bool

	// NewStyle returns a new Lip Gloss style that is suitable for the program's
	// environment.
	NewStyle() lipgloss.Style

	// ColorProfile returns the terminal's color profile.
	ColorProfile() lipgloss.Profile

	// what else?
}

type teaContext struct {
	context.Context

	profile         lipgloss.Profile
	kittyFlags      int
	backgroundColor color.Color
	hasLightBg      bool // cached value

	values map[interface{}]interface{}
	mtx    sync.Mutex
}

func newContext(ctx context.Context) *teaContext {
	c := new(teaContext)
	c.Context = ctx
	c.kittyFlags = -1
	c.values = make(map[interface{}]interface{})
	return c
}

func (c *teaContext) BackgroundColor() color.Color {
	return c.backgroundColor
}

func (c *teaContext) HasLightBackground() bool {
	return c.hasLightBg
}

func (c *teaContext) SupportsEnhancedKeyboard() bool {
	return c.kittyFlags >= 0
}

func (c *teaContext) NewStyle() lipgloss.Style {
	return lipgloss.NewStyle().ColorProfile(c.profile).HasLightBackground(c.hasLightBg)
}

func (c *teaContext) ColorProfile() lipgloss.Profile {
	return c.profile
}

func (ctx *teaContext) Value(key interface{}) interface{} {
	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()
	if v, ok := ctx.values[key]; ok {
		return v
	}
	return ctx.Context.Value(key)
}

func (ctx *teaContext) SetValue(key, value interface{}) {
	ctx.mtx.Lock()
	defer ctx.mtx.Unlock()
	ctx.values[key] = value
}
