package modules

const premiumSpoofCSS = `
[data-testid="upgrade-button"],
[data-testid="premium-link"],
[data-testid="upgrade-to-premium"],
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
[data-testid="consent-banner"] {
	display: none !important;
	visibility: hidden !important;
	opacity: 0 !important;
	pointer-events: none !important;
	height: 0 !important;
	width: 0 !important;
	overflow: hidden !important;
}

.main-nowPlayingBar-commercial {
	display: none !important;
}

[class*="sponsored"],
[class*="Sponsored"] {
	display: none !important;
}
`

const premiumSpoofJS = `
(function() {
if (window.__spotilitePremiumSpoofInstalled) return;
window.__spotilitePremiumSpoofInstalled = true;

function hidePremiumElements() {
var selectors = [
'[data-testid="upgrade-button"]',
'[data-testid="premium-link"]',
'[data-testid="upgrade-to-premium"]',
'button[aria-label*="Upgrade to Premium"]',
'button[aria-label*="Actualizar a Premium"]',
'button[aria-label*="Get Premium"]',
'button[aria-label*="Obtén Premium"]',
'[class*="upgrade-button"]',
'[class*="CookiePolicy"]',
'[class*="cookie-banner"]',
'[data-testid="cookie-banner"]'
];
selectors.forEach(function(sel) {
try {
var nodes = document.querySelectorAll(sel);
nodes.forEach(function(node) {
node.style.display = 'none';
node.style.visibility = 'hidden';
node.style.height = '0';
node.style.overflow = 'hidden';
});
} catch(e) {}
});

var iframes = document.querySelectorAll('iframe[src*="doubleclick"], iframe[src*="moatads"], iframe[src*="ads"]');
iframes.forEach(function(iframe) {
iframe.style.display = 'none';
iframe.style.height = '0';
});
}

setInterval(hidePremiumElements, 3000);
setTimeout(hidePremiumElements, 2000);
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
