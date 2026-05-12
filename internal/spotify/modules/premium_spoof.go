package modules

const premiumSpoofCSS = `
[data-testid="upgrade-button"],
[data-testid="premium-link"],
[data-testid="upgrade-to-premium"],
[data-testid="upsell-banner"],
[data-testid="premium-upsell"],
[href*="/premium"]:not([href*="/playlist"]):not([href*="/album"]):not([href*="/track"]):not([href*="/artist"]),
[href*="/subscription"],
.main-UpgradeButton,
.main-topBar-UpgradeButton,
.main-viewSelector-UpgradeButton,
button[aria-label="Upgrade to Premium"],
button[aria-label="Actualizar a Premium"],
button[aria-label="Get Premium"],
button[aria-label="Obtén Premium"],
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
'button[aria-label="Upgrade to Premium"]',
'button[aria-label="Actualizar a Premium"]',
'button[aria-label="Get Premium"]',
'button[aria-label="Obtén Premium"]',
'.main-UpgradeButton',
'.main-topBar-UpgradeButton',
'[class*="CookiePolicy"]',
'[class*="cookie-banner"]',
'[data-testid="cookie-banner"]',
'#onetrust-banner-sdk',
'#onetrust-consent-sdk'
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
var links = document.querySelectorAll('a[href*="/premium"], a[href*="/subscription"]');
for (var k = 0; k < links.length; k++) {
var href = links[k].href || '';
if (href.indexOf('/playlist/') !== -1 || href.indexOf('/album/') !== -1 || href.indexOf('/track/') !== -1 || href.indexOf('/artist/') !== -1) continue;
links[k].style.display = 'none';
}
} catch(e) {}

try {
var btns = document.querySelectorAll('button');
for (var m = 0; m < btns.length; m++) {
if (hasPremiumText(btns[m])) {
btns[m].style.display = 'none';
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

func (m *PremiumSpoofModule) CSS() string { return premiumSpoofCSS }
func (m *PremiumSpoofModule) JS() string  { return premiumSpoofJS }
func (m *PremiumSpoofModule) Selectors() []string {
	return []string{
		`[data-testid="upgrade-button"]`,
		`[data-testid="premium-link"]`,
		`[data-testid="upgrade-to-premium"]`,
		`button[aria-label="Upgrade to Premium"]`,
		`button[aria-label="Actualizar a Premium"]`,
		`button[aria-label="Get Premium"]`,
		`[class*="CookiePolicy"]`,
		`[class*="cookie-banner"]`,
	}
}
