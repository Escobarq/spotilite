package config

import "spotilite/internal/ini"

// parseINI dispatches to internal/ini.Parse, normalizing the section type so
// callers within the config package can keep operating on the local `sections`
// type alias without leaking that detail to consumers.
func parseINI(text string, out sections) error {
	shared := ini.Sections(out)
	return ini.Parse(text, shared)
}

func (c *Config) serialize() string {
	return ini.Serialize(ini.Sections(c.raw))
}
