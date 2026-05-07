package modules

type CustomCSSModule struct {
	BaseModule
}

func NewCustomCSSModule(enabled bool) *CustomCSSModule {
	return &CustomCSSModule{
		BaseModule: BaseModule{name: "custom_css", enabled: enabled},
	}
}

func (m *CustomCSSModule) CSS() string {
	if !m.enabled {
		return ""
	}
	return ""
}

func (m *CustomCSSModule) JS() string {
	if !m.enabled {
		return ""
	}
	return `(function() {
	var css = localStorage.getItem('spotilite.custom_css');
	if (!css) return;
	var id = 'spotilite-custom-css';
	var style = document.getElementById(id);
	if (!style) {
		style = document.createElement('style');
		style.id = id;
		document.head.appendChild(style);
	}
	style.textContent = css;
})();`
}
