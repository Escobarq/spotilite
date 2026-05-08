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
iframe[src*="doubleclick"],
iframe[src*="googleads"],
iframe[src*="googlesyndication"],
.WCbmOh4S3HVpA8RhH5Nj,
.Vs2HPUVcMf1MUfOb8KqE,
.tGKwoPuvNBNK3TzCSFQf,
.main-leaderboardComponent-container,
.main-billboardComponent-container,
[aria-label="Advertisement"],
[aria-label="Publicidad"],
[data-testid="card-clickout-ad"],
[data-testid="ad-card"],
[class*="ad-card"],
[class*="AdCard"],
[class*="ad-container"],
[class*="AdContainer"],
[class*="sponsor-card"],
[class*="SponsoredCard"],
[data-testid*="ad-"],
[class*="-ad-"],
[class*="_ad_"],
[class*="-Ad-"],
[class*="adElement"],
[class*="adElement"],
[class*="advertisement"],
[class*="Advertisement"],
[id*="ad-"],
[id*="Ad-"],
[id*="advertisement"],
[id*="sponsor"],
[class*="spotify-ad"],
[class*="SpotifyAd"],
[data-testid="inactive-ad"],
[class*="ad-break"],
[class*="AdBreak"],
[data-testid="ad-slot"],
[class*="ad-slot"],
[class*="AdSlot"],
[class*="stuck-ad"],
[class*="StuckAd"] {
	display: none !important;
	visibility: hidden !important;
	opacity: 0 !important;
	pointer-events: none !important;
	height: 0 !important;
	width: 0 !important;
	overflow: hidden !important;
	position: absolute !important;
	left: -9999px !important;
	top: -9999px !important;
}

.main-nowPlayingBar-NowPlayingBar > div[class*="ad"],
.main-nowPlayingBar-NowPlayingBar [data-testid="now-playing-bar-ad"],
.main-nowPlayingBar-NowPlayingBar [class*="Ad"],
.main-nowPlayingBar-NowPlayingBar [class*="ad"] {
	display: none !important;
}

*|*:not(html):not(body) > [class*="ad-break"]:not([class*="adblock"]),
*|*:not(html):not(body) > [class*="AdBreak"]:not([class*="adblock"]) {
	display: none !important;
}
`

const adblockJS = `
(function() {
if (window.__spotiliteAdblockInstalled) return;
window.__spotiliteAdblockInstalled = true;
window.__spotiliteIsAdPlaying = false;

var LOCAL_API = 'localhost:8765';
var AD_CONFIRM_THRESHOLD = 3;
var adConfirmCount = 0;

var AD_URL_PATTERNS = [
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
	'doubleclick.com',
	'googlesyndication.com',
	'googleadservices.com',
	'moatads.com',
	'moatads.net',
	'adservice.google.com',
	'adnxs.com',
	'adsrvr.org',
	'adsymptotic.com',
	'adform.net',
	'criteo.com',
	'rubiconproject.com',
	'pubmatic.com',
	'openx.net',
	'casalemedia.com',
	'contextweb.com',
	'rubicon',
	'criteo',
	'adtech',
	'adnxs',
	'adsrvr',
	'app-measurement',
	'analytics',
	'sp-logger',
	'adlog',
	'adobedtm',
	'tracking',
	'metrics',
	'telemetry'
];

var AD_SELECTORS = [
	'[data-testid="ad-type-banner"]',
	'[data-testid="billboard-ad"]',
	'[data-testid="leaderboard-ad"]',
	'[data-testid="sponsorship-ad"]',
	'[data-testid="sponsored-playlist"]',
	'[data-testid="hpto-ad"]',
	'[data-testid="ad-card"]',
	'[data-testid="ad-slot"]',
	'[data-testid="now-playing-bar-ad"]',
	'[data-testid="inactive-ad"]',
	'[class*="BillboardAd"]',
	'[class*="LeaderboardAd"]',
	'[class*="SponsoredPlaylist"]',
	'[class*="AdCard"]',
	'[class*="AdContainer"]',
	'[class*="ad-card"]',
	'[class*="ad-container"]',
	'[class*="sponsor-card"]',
	'[class*="ad-break"]',
	'[class*="AdBreak"]',
	'[class*="ad-slot"]',
	'[class*="AdSlot"]',
	'[class*="advertisement"]',
	'[aria-label="Advertisement"]',
	'[aria-label="Publicidad"]',
	'iframe[src*="doubleclick"]',
	'iframe[src*="moatads"]',
	'iframe[src*="googleads"]',
	'iframe[src*="ads"]',
	'iframe[src*="adsrvr"]',
	'iframe[src*="adform"]'
];

function isAdUrl(url) {
	if (!url) return false;
	var u = typeof url === 'string' ? url : (url.url ? url.url : '');
	if (!u || u.indexOf(LOCAL_API) !== -1) return false;
	for (var i = 0; i < AD_URL_PATTERNS.length; i++) {
		if (u.indexOf(AD_URL_PATTERNS[i]) !== -1) return true;
	}
	return false;
}

function createEmptyResponse() {
	return new Response('{}', {
		status: 200,
		headers: { 'Content-Type': 'application/json' }
	});
}

if (!window.__origFetch) {
	window.__origFetch = window.fetch.bind(window);
}

window.fetch = function() {
	var args = arguments;
	var url = args[0];

	if (isAdUrl(url)) {
		return Promise.resolve(createEmptyResponse());
	}

	var options = args[1];
	if (options && options.body) {
		try {
			var body = typeof options.body === 'string'
				? JSON.parse(options.body)
				: options.body;
			if (body && (body.type === 'ad' || body.typ === 'ad' || body.ad || body.isAd)) {
				return Promise.resolve(createEmptyResponse());
			}
		} catch(e) {}
	}

	return window.__origFetch.apply(this, args);
};

var _origXHROpen = XMLHttpRequest.prototype.open;
var _origXHRSend = XMLHttpRequest.prototype.send;

XMLHttpRequest.prototype.open = function(method, url) {
	if (isAdUrl(url)) {
		this.__adBlocked = true;
		return;
	}
	return _origXHROpen.apply(this, arguments);
};

XMLHttpRequest.prototype.send = function(data) {
	if (this.__adBlocked) {
		this.__adBlocked = false;
		return;
	}
	return _origXHRSend.apply(this, arguments);
};

function hideAdElements() {
	try {
		var hidden = 0;
		AD_SELECTORS.forEach(function(sel) {
			try {
				var els = document.querySelectorAll(sel);
				els.forEach(function(el) {
					if (el.style.display !== 'none') {
						el.style.setProperty('display', 'none', 'important');
						el.style.setProperty('visibility', 'hidden', 'important');
						el.style.setProperty('opacity', '0', 'important');
						el.style.setProperty('pointer-events', 'none', 'important');
						el.style.setProperty('height', '0', 'important');
						el.style.setProperty('width', '0', 'important');
						el.style.setProperty('position', 'absolute', 'important');
						el.style.setProperty('left', '-9999px', 'important');
						el.style.setProperty('top', '-9999px', 'important');
						hidden++;
					}
				});
			} catch(e) {}
		});

		if (window.__spotiliteDebug && hidden > 0) {
			console.log('[Spotilite AdBlock] Hidden ' + hidden + ' elements');
		}
	} catch(e) {}
}

function trySkipAd() {
	if (!window.__spotiliteIsAdPlaying) {
		var isAd = checkIfAdPlaying();
		if (!isAd) return false;
		adConfirmCount++;
		if (adConfirmCount < AD_CONFIRM_THRESHOLD) {
			console.log('[Spotilite AdBlock] Ad detected, waiting for confirmation (' + adConfirmCount + '/' + AD_CONFIRM_THRESHOLD + ')');
			return false;
		}
		console.log('[Spotilite AdBlock] Ad confirmed, attempting skip');
		window.__spotiliteIsAdPlaying = true;
	}

	try {
		var skipBtn = document.querySelector('[data-testid="skip-button"]');
		if (skipBtn && skipBtn.offsetParent !== null && !skipBtn.disabled) {
			skipBtn.click();
			console.log('[Spotilite AdBlock] Skipped ad');
			resetAdState();
			return true;
		}

		var nextBtn = document.querySelector('[data-testid="control-button-skip-forward"]');
		if (nextBtn && nextBtn.offsetParent !== null && !nextBtn.disabled) {
			nextBtn.click();
			console.log('[Spotilite AdBlock] Skipped ad via next');
			resetAdState();
			return true;
		}
	} catch(e) {}
	return false;
}

function checkIfAdPlaying() {
	try {
		var nowPlayingBar = document.querySelector('.main-nowPlayingBar-nowPlayingBar');
		if (!nowPlayingBar) return false;

		var adBadge = nowPlayingBar.querySelector('[data-testid="now-playing-bar-ad"]');
		if (adBadge) return true;

		var entityBadges = nowPlayingBar.querySelectorAll('[class*="badge"], [class*="Badge"], [data-testid*="ad"]');
		for (var i = 0; i < entityBadges.length; i++) {
			var text = entityBadges[i].textContent || '';
			if (text.toLowerCase().includes('ad') || text.toLowerCase().includes('sponsor') || text.toLowerCase().includes('publicidad')) {
				return true;
			}
		}

		var skipBtn = document.querySelector('[data-testid="skip-button"]');
		if (skipBtn) {
			var ariaLabel = skipBtn.getAttribute('aria-label') || '';
			if (ariaLabel.toLowerCase().includes('skip')) {
				var isVisible = skipBtn.offsetParent !== null && !skipBtn.disabled;
				if (isVisible) return true;
			}
		}

		var titleContainer = nowPlayingBar.querySelector('[class*="title"], [class*="Title"]');
		if (titleContainer) {
			var titleText = titleContainer.textContent || '';
			if (titleText.includes(':')) {
				var parts = titleText.split(':');
				if (parts.length === 2 && parts[0].trim().length < 20 && parts[1].trim().length < 20) {
					return true;
				}
			}
		}
	} catch(e) {}
	return false;
}

function resetAdState() {
	adConfirmCount = 0;
	window.__spotiliteIsAdPlaying = false;
}

function detectAndBlockAudioAds() {
	var isAd = checkIfAdPlaying();
	if (isAd && adConfirmCount >= AD_CONFIRM_THRESHOLD) {
		trySkipAd();
	}
}

function detectAndBlockAudioAds() {
	var isAd = checkIfAdPlaying();
	if (isAd && adConfirmCount >= AD_CONFIRM_THRESHOLD) {
		trySkipAd();
	}
}

hideAdElements();

var observer = new MutationObserver(function(mutations) {
	var shouldHide = false;
	var shouldCheckSkip = false;

	mutations.forEach(function(mutation) {
		if (mutation.type === 'childList') {
			mutation.addedNodes.forEach(function(node) {
				if (node.nodeType === Node.ELEMENT_NODE) {
					AD_SELECTORS.forEach(function(sel) {
						try {
							if (node.matches && node.matches(sel)) {
								shouldHide = true;
							}
							var children = node.querySelectorAll(sel);
							if (children.length > 0) shouldHide = true;
						} catch(e) {}
					});
				}
			});
		}

		if (mutation.type === 'attributes') {
			var attr = mutation.attributeName;
			if (attr === 'class' || attr === 'data-testid' || attr === 'aria-label') {
				var val = mutation.target.getAttribute(attr) || '';
				if (val.toLowerCase().includes('ad') || val.toLowerCase().includes('sponsor')) {
					shouldHide = true;
					shouldCheckSkip = true;
				}
			}
		}
	});

	if (shouldHide) hideAdElements();
	if (shouldCheckSkip) trySkipAd();
});

observer.observe(document.body || document.documentElement, {
	childList: true,
	subtree: true,
	attributes: true,
	attributeFilter: ['class', 'data-testid', 'aria-label', 'src', 'href']
});

var checkCount = 0;
var skipCooldown = 0;
var checkInterval = setInterval(function() {
	checkCount++;
	hideAdElements();

	if (skipCooldown > 0) {
		skipCooldown--;
	}

	if (checkCount % 3 === 0 && skipCooldown === 0) {
		var wasAd = window.__spotiliteIsAdPlaying;
		var skipped = trySkipAd();
		if (skipped) {
			skipCooldown = 6;
		}
		detectAndBlockAudioAds();
	}

	if (checkCount > 1000) {
		checkCount = 0;
	}

	if (document.readyState === 'complete' && checkCount % 5 === 0) {
		var adOverlay = document.querySelector('.ad-overlay, .ad-backdrop, [class*="ad-overlay"], [class*="AdOverlay"]');
		if (adOverlay && adOverlay.offsetParent !== null) {
			adOverlay.style.setProperty('display', 'none', 'important');
		}
	}

	if (checkCount % 10 === 0 && !window.__spotiliteIsAdPlaying) {
		adConfirmCount = 0;
	}
}, 500);

window.__hideAdElements = hideAdElements;
window.__trySkipAd = trySkipAd;
window.__checkIfAdPlaying = checkIfAdPlaying;

console.log('[Spotilite AdBlock] Initialized with confirmation-based skip (' + AD_CONFIRM_THRESHOLD + ' confirms required)');
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
		`[class*="ad-card"]`,
		`[class*="AdCard"]`,
		`[class*="ad-container"]`,
		`[class*="AdContainer"]`,
		`[class*="ad-break"]`,
		`[class*="AdBreak"]`,
		`[class*="ad-slot"]`,
		`[class*="AdSlot"]`,
	}
}