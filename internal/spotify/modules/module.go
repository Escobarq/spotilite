package modules

type Module interface {
	Name() string
	CSS() string
	JS() string
	Selectors() []string
	Enabled() bool
	SetEnabled(bool)
}

type BaseModule struct {
	name    string
	enabled bool
}

func (b *BaseModule) Name() string    { return b.name }
func (b *BaseModule) Enabled() bool   { return b.enabled }
func (b *BaseModule) SetEnabled(v bool) { b.enabled = v }
func (b *BaseModule) CSS() string     { return "" }
func (b *BaseModule) JS() string      { return "" }
func (b *BaseModule) Selectors() []string { return nil }
