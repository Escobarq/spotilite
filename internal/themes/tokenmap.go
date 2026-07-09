package themes

// TokenMap bridges Spicetify's canonical --spice-* variable names to the
// Spotify web player's own design tokens (the CSS variables actually consumed
// by open.spotify.com's bundled CSS).
//
// Multiple web tokens can map to one spice variable; they are all aliased to
// the same source. Themes authored for spicetify-cli that only know about
// --spice-* names continue to work, and the visual effect propagates to the
// web player's actual styling.
//
// The list was derived by inspecting open.spotify.com's compiled CSS for the
// most commonly referenced custom properties at :root scope. It is
// intentionally conservative: unmapped --spice-* variables are still emitted
// (themes can target them directly), and web tokens not present in this map
// are untouched so Spotify's defaults remain.
//
// Spicetify variable names follow:
// https://spicetify.app/docs/development/themes (color.ini reference).
var TokenMap = map[string][]string{
	"main": {
		"--background-base",
		"--background-noise-base",
	},
	"main-elevated": {
		"--background-elevated-base",
	},
	"highlight": {
		"--background-press",
	},
	"highlight-elevated": {
		"--background-elevated-highlight",
	},
	"sidebar": {
		"--background-tinted-base",
	},
	"player": {
		"--background-tinted-base",
	},
	"card": {
		"--background-elevated-base",
	},
	"card-elevated": {
		"--background-press",
	},
	"button": {
		"--essential-color-accent",
		"--essential-bright-accent",
	},
	"button-active": {
		"--essential-color-accent",
	},
	"button-disabled": {
		"--essential-subdued",
	},
	"tab-active": {
		"--background-tinted-base",
		"--background-elevated-highlight",
	},
	"text": {
		"--text-base",
		"--text-subdued",
	},
	"subtext": {
		"--text-subdued",
		"--text-bright-accent",
	},
	"notification": {
		"--background-elevated-highlight",
	},
	"notification-error": {
		"--essential-negative",
	},
	"misc": {
		"--essential-subdued",
	},
}
