package modules

const lyricsRulesCSS = `
.Li269NgzkU2gI4KOP9sM,
.I2WIloMMjsBeMaIS8H3v,
.McI3hD7aCfpq015LJa6X,
.gpDSOimnzH4zTJmE7UR5 {
	--lyrics-color-active: var(--lyrics-current) !important;
	--lyrics-color-inactive: var(--lyrics-next) !important;
	--lyrics-color-passed: var(--lyrics-past) !important;
	--lyrics-color-background: var(--lyrics-bg) !important;
}

p[class*="e-"][class*="-text"].encore-text-body-small {
	color: var(--lyrics-mm) !important;
	margin-bottom: 8px !important;
}

[data-testid="lyrics-npv-section"]:not(._OhUGn8Plh3mRw4awIM5):not(.Sb2rC16jDkGc9eweOU8g):not(._YRfjT5prbRuSXcNK9WR):not(.RXRGSIFllAhUYWKYlANd) {
	background-color: #1f1f1f !important;
}

.ebHsEf.I4K12o0qDoITOLr2AEs0,
.ebHsEf.OYiGFGZJDIZ4FF4ZTDK2,
.jKdLzW.LvLs_UgYs7ps5KdoCr0h,
.bWzOVV._T5UDP2tItG9WGdwO5Yi,
.hzUuLPdH48AzgQun5NYQ [data-encore-id="type"],
.hzUuLPdH48AzgQun5NYQ [data-encore-id="text"],
body .LomBcMvfM8AEmZGquAdj,
body .W_EplVEAbZrZURqfLiQC,
body .kGR_hu4tdj9PnUlSPaRL,
.GML6YUVCeJvRhGznLnqm,
body .iq4cgi0YEKr6DGaTtzUj,
body .KDhLFoEqoClhH12bsfrS {
	color: var(--lyrics-mm) !important;
}

.C3pBU1DsOUJJOAv89ZFT,
.T67LFP0PElpfkkLuegQt,
.e7eFLioNSG5PAi1qVFT4,
.l2060qoyWU4J9ihSxHLE,
.hfTlyhd7WCIk9xmP {
	color: var(--lyrics-mm) !important;
}

.FUYNhisXTCmbzt9IDxnT,
.tr8V5eHsUaIkOYVw7eSG,
.hW9km7ku6_iggdWDR_Lg,
.lofIAg8Ixko3mfBrbfej,
.bbJIIopLxggQmv5x {
	--lyrics-color-active: var(--lyrics-current) !important;
	--lyrics-color-inactive: var(--lyrics-next) !important;
	--lyrics-color-passed: var(--lyrics-past) !important;
	--lyrics-color-background: var(--lyrics-bg) !important;
}

.H2J92dVdr0ykdOX5azL1,
.KnFq2ijXFdOtyl4Iebjv {
	color: var(--lyrics-past) !important;
	opacity: 1 !important;
}

.vapgYYF2HMEeLJuOWGq5:hover,
._LKG3z7SnerR0eigPCoK:hover,
.NHVfxGs2HwmI_fly2JC4:hover,
.gaHIufRWhoWbiT8S6zuM:hover,
.FQYXZaa0aDIrse54YlYO:hover,
.PSyA4iign083ZV6vOqPj:hover,
.Plu0zvuRv7kOQwsQ02cC:hover,
.Mnf9PkrVHsX90BNf:hover,
.XiH9KR6bhDwEFykV:hover {
	color: var(--lyrics-hover) !important;
	text-decoration: none !important;
}

.HxblHEsl2WX2yhubfVIc,
.SruqsAzX8rUtY2isUZDF,
.eTLjCqbDo7QehPEPz86a,
.AEfhRyqGa3vzQrgfdwWE.Re403AJffPPuZmX7LRJj,
.NHVfxGs2HwmI_fly2JC4.E64X_eoy6xsJmDdKKHja,
.gaHIufRWhoWbiT8S6zuM.Qo3OkrSis5IWlP9Tchbr,
.AEfhRyqGa3vzQrgfdwWE .Re403AJffPPuZmX7LRJj {
	color: var(--lyrics-next) !important;
}

.npv-lyrics__text-wrapper--previous .npv-lyrics__text {
	color: var(--lyrics-past) !important;
}
.npv-lyrics__text-wrapper--current .npv-lyrics__text {
	color: var(--lyrics-current) !important;
}
.npv-lyrics__text-wrapper--next .npv-lyrics__text {
	color: var(--lyrics-next) !important;
}
.npv-lyrics__text.npv-lyrics__text--credits,
.npv-lyrics__text--unsynced-warning {
	color: var(--lyrics-mm) !important;
}
.npv-lyrics__text--unsynced {
	color: var(--lyrics-next) !important;
}
.npv-background-color {
	background: var(--lyrics-bg) !important;
}
.npv-main-container {
	background: transparent !important;
}
.npv-lyrics__gradient-background {
	background: -webkit-gradient(linear, left top, left bottom, from(rgba(18,18,18,0)), color-stop(30%, var(--lyrics-bg)), color-stop(60%, var(--lyrics-bg))) !important;
	background: -webkit-linear-gradient(top, rgba(18,18,18,0) 0%, var(--lyrics-bg) 30%, var(--lyrics-bg) 60%) !important;
	background: linear-gradient(to bottom, rgba(18,18,18,0) 0%, var(--lyrics-bg) 30%, var(--lyrics-bg) 60%) !important;
}

.l6lFMYQteTVnTcHnLywc,
._nDkCIVgkWayq3tqiIuW,
.B_wut2Bw4HwLr3w8rNfM {
	--transcript-color-background: var(--lyrics-bg) !important;
	--transcript-color-text: var(--lyrics-next) !important;
	--transcript-color-highlightText: var(--lyrics-current) !important;
}
`

type LyricsThemeModule struct {
	BaseModule
	theme string
}

type ThemeColors struct {
	Past  string `json:"past"`
	Next  string `json:"next"`
	Hover string `json:"hover"`
	Bg    string `json:"bg"`
	MM    string `json:"mm"`
}

var LyricsThemes = map[string]ThemeColors{
	"default":     {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"red":         {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"orange":      {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"yellow":      {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"spotify":     {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"spotify-alt": {Past: "#9b9b9b", Next: "#666666", Hover: "#f2f2f2", Bg: "#242424", MM: "#C2C2C2"},
	"blue":        {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"purple":      {Past: "#575757", Next: "#575757", Hover: "#C8C8C8", Bg: "#121212", MM: "#969696"},
	"strawberry":  {Past: "#F17F7F", Next: "#595959", Hover: "#F2F2F2", Bg: "#1C1C1E", MM: "#595959"},
	"pumpkin":     {Past: "#FDAC69", Next: "#595959", Hover: "#F2F2F2", Bg: "#1C1C1E", MM: "#595959"},
	"sandbar":     {Past: "#FFDB7A", Next: "#595959", Hover: "#F2F2F2", Bg: "#1C1C1E", MM: "#595959"},
	"radium":      {Past: "#AAFFA3", Next: "#595959", Hover: "#F2F2F2", Bg: "#1C1C1E", MM: "#595959"},
	"oceano":      {Past: "#70DBF0", Next: "#595959", Hover: "#F2F2F2", Bg: "#1C1C1E", MM: "#595959"},
	"royal":       {Past: "#B8A3EB", Next: "#595959", Hover: "#F2F2F2", Bg: "#1C1C1E", MM: "#595959"},
	"github":      {Past: "#AD82F8", Next: "#47566D", Hover: "#70B3FF", Bg: "#161B22", MM: "#408BD0"},
	"discord":     {Past: "#616774", Next: "#616774", Hover: "#FFFFFF", Bg: "#23272A", MM: "#616774"},
	"drot":        {Past: "#505050", Next: "#505050", Hover: "#A13131", Bg: "#191414", MM: "#787878"},
	"forest":      {Past: "#505050", Next: "#505050", Hover: "#418022", Bg: "#141914", MM: "#787878"},
	"fresh":       {Past: "#505050", Next: "#505050", Hover: "#0B7383", Bg: "#14191E", MM: "#787878"},
	"zing":        {Past: "#4E596F", Next: "#4E596F", Hover: "#FFFFFF", Bg: "#202430", MM: "#9EA8BC"},
	"pinkle":      {Past: "#9579E3", Next: "#5E547C", Hover: "#FFFFFF", Bg: "#1C1925", MM: "#5E547C"},
	"krux":        {Past: "#5C89D2", Next: "#696E79", Hover: "#FFFFFF", Bg: "#191E29", MM: "#696E79"},
	"blueberry":   {Past: "#1CAAC6", Next: "#516377", Hover: "#A0D1FA", Bg: "#232937", MM: "#516377"},
	"postlight":   {Past: "#C9A8FE", Next: "#534D6F", Hover: "#D1D1D1", Bg: "#13101C", MM: "#534D6F"},
	"relish":      {Past: "#9D2117", Next: "#C8A032", Hover: "#E5CB8B", Bg: "#121212", MM: "#787878"},
	"turquoise":   {Past: "#00656aa0", Next: "#575757", Hover: "#a97aff", Bg: "#121212", MM: "#00656a"},
	"lavender":    {Past: "#B8A2EA", Next: "#575757", Hover: "#F2F2F2", Bg: "#121212", MM: "#C2C2C2"},
}

var themeCurrentColors = map[string]string{
	"default":     "#C8C8C8",
	"red":         "#FF3737",
	"orange":      "#F68D00",
	"yellow":      "#ECE224",
	"spotify":     "#1ED760",
	"spotify-alt": "#1ed760",
	"blue":        "#00DFEA",
	"purple":      "#9E6BE3",
	"strawberry":  "#E43A47",
	"pumpkin":     "#E88500",
	"sandbar":     "#F5BA18",
	"radium":      "#17D344",
	"oceano":      "#13A1BD",
	"royal":       "#8461DD",
	"github":      "#7EE787",
	"discord":     "#7A8FDC",
	"drot":        "#F37171",
	"forest":      "#AEF97B",
	"fresh":       "#50DCF0",
	"zing":        "#F67064",
	"pinkle":      "#CD3B99",
	"krux":        "#01C38D",
	"blueberry":   "#90E0F0",
	"postlight":   "#9D65C7",
	"relish":      "#C8C8C8",
	"turquoise":   "#01dfea",
	"lavender":    "#8462DD",
}

var DefaultTheme = "spotify"

func NewLyricsThemeModule(enabled bool, theme string) *LyricsThemeModule {
	m := &LyricsThemeModule{
		BaseModule: BaseModule{name: "lyrics_theme", enabled: enabled},
	}
	m.SetTheme(theme)
	return m
}

func (m *LyricsThemeModule) SetTheme(theme string) {
	if _, ok := LyricsThemes[theme]; !ok {
		theme = DefaultTheme
	}
	m.theme = theme
}

func (m *LyricsThemeModule) CSS() string {
	if !m.enabled {
		return ""
	}
	c, ok := LyricsThemes[m.theme]
	if !ok {
		c = LyricsThemes[DefaultTheme]
	}
	cur, ok := themeCurrentColors[m.theme]
	if !ok {
		cur = themeCurrentColors[DefaultTheme]
	}
	return `:root {
--lyrics-past: ` + c.Past + `;
--lyrics-current: ` + cur + `;
--lyrics-next: ` + c.Next + `;
--lyrics-hover: ` + c.Hover + `;
--lyrics-bg: ` + c.Bg + `;
--lyrics-mm: ` + c.MM + `;
}
` + lyricsRulesCSS
}

func (m *LyricsThemeModule) Theme() string {
	return m.theme
}

func ThemeList() []string {
	themes := make([]string, 0, len(LyricsThemes))
	for k := range LyricsThemes {
		themes = append(themes, k)
	}
	return themes
}

func ThemeColorsList() map[string]string {
	colors := make(map[string]string)
	for name := range LyricsThemes {
		if c, ok := themeCurrentColors[name]; ok {
			colors[name] = c
		}
	}
	return colors
}
