package modules

const adblockCSS = `
[data-testid="ad-type-banner"],
[data-testid="billboard-ad"],
[data-testid="leaderboard-ad"],
[data-testid="sponsorship-ad"],
[data-testid="sponsored-playlist"],
[data-testid="hpto-ad"],
[class*="ad-banner"],
[class*="BillboardAd"],
[class*="LeaderboardAd"],
[class*="SponsoredPlaylist"],
[class*="hpto-"],
iframe[src*="ad.doubleclick.net"],
iframe[src*="moatads"],
iframe[src*="ads"],
.WCbmOh4S3HVpA8RhH5Nj,
.Vs2HPUVcMf1MUfOb8KqE,
.tGKwoPuvNBNK3TzCSFQf,
.main-leaderboardComponent-container,
.main-billboardComponent-container,
[aria-label="Advertisement"],
[aria-label="Publicidad"],
[data-testid="card-clickout-ad"],
[data-testid="ad-card"] {
	display: none !important;
	visibility: hidden !important;
	opacity: 0 !important;
	pointer-events: none !important;
	height: 0 !important;
	width: 0 !important;
	overflow: hidden !important;
}

.main-nowPlayingBar-NowPlayingBar > div[class*="ad"],
.main-nowPlayingBar-NowPlayingBar [data-testid="now-playing-bar-ad"] {
	display: none !important;
}
`

const adblockJS = `
(function() {
if (window.__spotiliteAdblockInstalled) return;
window.__spotiliteAdblockInstalled = true;

if (!window.__origFetch) window.__origFetch = window.fetch.bind(window);
var _origFetch = window.__origFetch;

var LOCAL_API = 'localhost:8765';
var AD_URLS = [
'ad-handler.spotify.com',
'ad-return-url',
'gaia.spotify.com/ad',
'spclient.wg.spotify.com/ad',
'partnerakamai.spotify.com/ad',
'/ad/',
'ads.php',
'spotify.com/ad-logic',
'spotify.com/ads/',
'spotify.com/pair/',
'spotify.com/gabo-receiver-service',
'doubleclick.net',
'googlesyndication',
'moatads',
'adservice.google',
'adnxs.com',
'adsrvr.org'
];

function isAdUrl(url) {
var u = typeof url === 'string' ? url : (url && url.url ? url.url : '');
if (!u) return false;
if (u.indexOf(LOCAL_API) !== -1) return false;
for (var i = 0; i < AD_URLS.length; i++) {
if (u.indexOf(AD_URLS[i]) !== -1) return true;
}
return false;
}

window.fetch = function() {
var args = arguments;
var url = args[0];
if (isAdUrl(url)) {
return Promise.resolve(new Response('{}', {
status: 200,
headers: { 'Content-Type': 'application/json' }
}));
}

var options = args[1];
if (options && options.body && typeof options.body === 'string') {
try {
var body = JSON.parse(options.body);
if (body && (body.typ === 'ad' || body.type === 'ad')) {
return Promise.resolve(new Response('{}', {
status: 200,
headers: { 'Content-Type': 'application/json' }
}));
}
} catch(e) {}
}

return _origFetch.apply(this, args);
};

var _origXHROpen = XMLHttpRequest.prototype.open;
XMLHttpRequest.prototype.open = function(method, url) {
if (isAdUrl(url)) {
return;
}
return _origXHROpen.apply(this, arguments);
};

function skipAds() {
var adIndicator = document.querySelector('[data-testid="ad-type-banner"], [class*="ad-indicator"], [data-testid="now-playing-bar-ad"]');
if (adIndicator) {
var skipBtn = document.querySelector('[data-testid="skip-button"], [class*="skip-ad"], button[aria-label*="Skip"]');
if (skipBtn) skipBtn.click();
}

var nowPlaying = document.querySelector('.main-nowPlayingBar-nowPlayingBar');
if (nowPlaying) {
var adBar = nowPlaying.querySelector('[class*="ad"], [data-testid*="ad"]');
if (adBar) {
var nextBtn = document.querySelector('[data-testid="control-button-skip-forward"], button[aria-label="Next"], button[aria-label="Siguiente"]');
if (nextBtn) nextBtn.click();
}
}
}

setInterval(skipAds, 2000);
})();
`

type AdBlockModule struct {
	BaseModule
}

func NewAdBlockModule(enabled bool) *AdBlockModule {
	return &AdBlockModule{
		BaseModule: BaseModule{name: "adblock", enabled: enabled},
	}
}

func (m *AdBlockModule) CSS() string     { return adblockCSS }
func (m *AdBlockModule) JS() string      { return adblockJS }
func (m *AdBlockModule) Selectors() []string {
	return []string{
		`[data-testid="ad-type-banner"]`,
		`[data-testid="billboard-ad"]`,
		`[data-testid="leaderboard-ad"]`,
		`[data-testid="sponsorship-ad"]`,
		`[data-testid="sponsored-playlist"]`,
		`[data-testid="hpto-ad"]`,
		`[class*="BillboardAd"]`,
		`[class*="LeaderboardAd"]`,
		`[aria-label="Advertisement"]`,
		`[aria-label="Publicidad"]`,
	}
}
