package spicetify

import _ "embed"

// Helper scripts ported verbatim from spicetify-cli (cli-main/jsHelper).
// They are emitted into the bundle in the same order as desktop spicetify
// does (see cli-main/src/apply/apply.go htmlMod). Each one is a standalone
// IIFE that runs on page load and expects `window.Spicetify` to exist by the
// time it executes (core.js + shims run first).

//go:embed helpers/expFeatures.js
var expFeaturesJS string

//go:embed helpers/homeConfig.js
var homeConfigJS string

//go:embed helpers/sidebarConfig.js
var sidebarConfigJS string
