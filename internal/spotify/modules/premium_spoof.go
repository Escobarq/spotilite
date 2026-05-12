package modules

const premiumSpoofCSS = `
[data-testid="upgrade-button"],
[data-testid="premium-link"],
[data-testid="upgrade-to-premium"],
[data-testid="upsell-banner"],
[data-testid="premium-upsell"],
[href*="/premium"],
[href*="/subscription"],
[class*="upgrade-button"],
[class*="upsell"],
[class*="premium-upsell"],
.main-UpgradeButton,
.main-topBar-UpgradeButton,
.main-viewSelector-UpgradeButton,
button[aria-label*="Upgrade to Premium"],
button[aria-label*="Actualizar a Premium"],
button[aria-label*="Get Premium"],
button[aria-label*="Obtén Premium"],
.main-billboardComponent,
.main-leaderboardComponent,
.WCbmOh4S3HVpA8RhH5Nj,
.Vs2HPUVcMf1MUfOb8KqE,
[class*="CookiePolicy"],
[class*="cookie-banner"],
#onetrust-banner-sdk,
#onetrust-consent-sdk,
[data-testid="cookie-banner"],
[data-testid="consent-banner"],
[data-encore-id="buttonTertiary"] {
	display: none !important;
}
`

const premiumSpoofJS = `
(function() {
if (window.__spotilitePremiumSpoofInstalled) return;
window.__spotilitePremiumSpoofInstalled = true;

var selectors = [
	'[data-testid="upgrade-button"]',
	'[data-testid="premium-link"]',
	'[data-testid="upgrade-to-premium"]',
	'[data-testid="upsell-banner"]',
	'[data-testid="premium-upsell"]',
	'button[aria-label*="Upgrade to Premium"]',
	'button[aria-label*="Actualizar a Premium"]',
	'button[aria-label*="Get Premium"]',
	'button[aria-label*="Obtén Premium"]',
	'[class*="upgrade-button"]',
	'[class*="upsell"]',
	'[class*="premium-upsell"]',
	'[class*="CookiePolicy"]',
	'[class*="cookie-banner"]',
	'[data-testid="cookie-banner"]',
	'a[href*="/premium"]',
	'a[href*="/subscription"]',
	'[data-encore-id="buttonTertiary"]',
	'[data-encore-id="buttonSecondary"]'
];

var premiumTexts = ['descubrir premium', 'discover premium', 'get premium', 'upgrade to premium', 'actualizar a premium', 'obtén premium', 'go premium'];

function hasPremiumText(el) {
	var text = (el.textContent || '').toLowerCase().trim();
	for (var i = 0; i < premiumTexts.length; i++) {
		if (text.indexOf(premiumTexts[i]) !== -1) return true;
	}
	return false;
}

function hide() {
	for (var i = 0; i < selectors.length; i++) {
		try {
			var els = document.querySelectorAll(selectors[i]);
			for (var j = 0; j < els.length; j++) {
				els[j].style.display = 'none';
			}
		} catch(e) {}
	}

	try {
		var btns = document.querySelectorAll('button, a');
		for (var k = 0; k < btns.length; k++) {
			if (hasPremiumText(btns[k])) {
				btns[k].style.display = 'none';
			}
		}
	} catch(e) {}
}

hide();
setInterval(hide, 3000);
})();
`

type PremiumSpoofModule struct {
	BaseModule
}

func NewPremiumSpoofModule(enabled bool) *PremiumSpoofModule {
	return &PremiumSpoofModule{
		BaseModule: BaseModule{name: "premium_spoof", enabled: enabled},
	}
}

func (m *PremiumSpoofModule) CSS() string     { return premiumSpoofCSS }
func (m *PremiumSpoofModule) JS() string      { return premiumSpoofJS }
func (m *PremiumSpoofModule) Selectors() []string {
	return []string{
		`[data-testid="upgrade-button"]`,
		`[data-testid="premium-link"]`,
		`[data-testid="upgrade-to-premium"]`,
		`[href*="/premium"]`,
		`button[aria-label*="Upgrade to Premium"]`,
		`button[aria-label*="Actualizar a Premium"]`,
		`button[aria-label*="Get Premium"]`,
		`[class*="CookiePolicy"]`,
		`[class*="cookie-banner"]`,
	}
}
