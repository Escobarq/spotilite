package modules

const adblockCSS = `
/* === AD BLOCKING - SPECIFIC SELECTORS ONLY === */
[data-testid="ad-type-banner"],
[data-testid="billboard-ad"],
[data-testid="leaderboard-ad"],
[data-testid="sponsorship-ad"],
[data-testid="sponsored-playlist"],
[data-testid="hpto-ad"],
[data-testid="card-clickout-ad"],
[data-testid="ad-card"],
[data-testid="ad-slot"],
[data-testid="inactive-ad"],
[data-testid="now-playing-bar-ad"],
[class*="BillboardAd"],
[class*="LeaderboardAd"],
[class*="SponsoredPlaylist"],
[class*="ad-banner"],
[class*="ad-card"],
[class*="AdCard"],
[class*="sponsor-card"],
[class*="SponsoredCard"],
[class*="ad-slot"],
[class*="AdSlot"],
[class*="stuck-ad"],
[class*="StuckAd"],
[class*="spotify-ad"],
[class*="SpotifyAd"],
.main-leaderboardComponent-container,
.main-billboardComponent-container,
[aria-label="Advertisement"],
[aria-label="Publicidad"],
iframe[src*="ad.doubleclick.net"],
iframe[src*="moatads"],
iframe[src*="ads"],
iframe[src*="doubleclick"],
iframe[src*="googleads"],
iframe[src*="googlesyndication"],
.WCbmOh4S3HVpA8RhH5Nj,
.Vs2HPUVcMf1MUfOb8KqE,
.tGKwoPuvNBNK3TzCSFQf {
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

/* Ad break overlays and modals - more specific */
body > [class*="ad-break"],
body > [class*="AdBreak"],
div[class*="ad-break"] > *,
div[class*="AdBreak"] > * {
	display: none !important;
	visibility: hidden !important;
	opacity: 0 !important;
	pointer-events: none !important;
}

/* === CRITICAL: PROTECT NOW PLAYING BAR AND ALL CHILDREN === */
.main-nowPlayingBar-nowPlayingBar,
.main-nowPlayingBar-NowPlayingBar,
[data-testid="now-playing-bar"],
.now-playing-bar,
footer[class*="NowPlayingBar"],
footer[class*="nowPlayingBar"] {
	display: flex !important;
	visibility: visible !important;
	opacity: 1 !important;
	pointer-events: auto !important;
	position: relative !important;
	height: auto !important;
	min-height: 90px !important;
	z-index: 1000 !important;
}

/* Protect ALL direct children and descendants of now playing bar */
.main-nowPlayingBar-nowPlayingBar *,
.main-nowPlayingBar-NowPlayingBar *,
[data-testid="now-playing-bar"] *,
.now-playing-bar *,
footer[class*="NowPlayingBar"] *,
footer[class*="nowPlayingBar"] * {
	display: inherit !important;
	visibility: visible !important;
	opacity: 1 !important;
	pointer-events: auto !important;
	position: static !important;
	max-height: none !important;
	min-height: 0 !important;
}

/* Specifically protect track info area */
.main-nowPlayingBar-nowPlayingBar [data-testid="track-info"],
.main-nowPlayingBar-nowPlayingBar .track-info,
.main-nowPlayingBar-nowPlayingBar [class*="TrackInfo"],
[data-testid="now-playing-bar"] [data-testid="track-info"],
[data-testid="now-playing-bar"] .track-info,
[data-testid="now-playing-bar"] [class*="TrackInfo"] {
	display: flex !important;
	visibility: visible !important;
	opacity: 1 !important;
	pointer-events: auto !important;
	position: static !important;
	flex-direction: row !important;
	align-items: center !important;
}

/* Protect playback controls */
.main-nowPlayingBar-nowPlayingBar [data-testid="control-button"],
.main-nowPlayingBar-nowPlayingBar [class*="PlaybackControl"],
.main-nowPlayingBar-nowPlayingBar [class*="PlayerControls"],
[data-testid="now-playing-bar"] [data-testid="control-button"],
[data-testid="now-playing-bar"] [class*="PlaybackControl"],
[data-testid="now-playing-bar"] [class*="PlayerControls"] {
	display: inline-flex !important;
	visibility: visible !important;
	opacity: 1 !important;
	pointer-events: auto !important;
	position: static !important;
}

/* Protect volume and extra controls */
.main-nowPlayingBar-nowPlayingBar [data-testid="volume-bar"],
.main-nowPlayingBar-nowPlayingBar [class*="VolumeBar"],
.main-nowPlayingBar-nowPlayingBar [data-testid="extra-controls"],
.main-nowPlayingBar-nowPlayingBar [class*="ExtraControls"] {
	display: flex !important;
	visibility: visible !important;
	opacity: 1 !important;
	pointer-events: auto !important;
	position: static !important;
}

/* Override any ad-related styles within now playing bar */
.main-nowPlayingBar-nowPlayingBar [class*="Commercial"],
.main-nowPlayingBar-nowPlayingBar [class*="AdBreak"],
.main-nowPlayingBar-nowPlayingBar [data-testid="ad-"],
[data-testid="now-playing-bar"] [class*="Commercial"],
[data-testid="now-playing-bar"] [class*="AdBreak"],
[data-testid="now-playing-bar"] [data-testid="ad-"] {
	display: none !important;
	visibility: hidden !important;
	opacity: 0 !important;
	pointer-events: none !important;
	height: 0 !important;
	width: 0 !important;
	padding: 0 !important;
	margin: 0 !important;
}
`

const adblockJS = `
(function() {
if (window.__spotiliteAdblockInstalled) return;
window.__spotiliteAdblockInstalled = true;
window.__spotiliteIsAdPlaying = false;
window.__spotiliteAdConfirmCount = 0;
window.__spotiliteLastCheck = 0;
window.__spotiliteIdleTime = 0;

var LOCAL_API = 'localhost:8765';
var AD_CONFIRM_THRESHOLD = 2;
var CHECK_INTERVAL_BASE = 1000;
var CHECK_INTERVAL_IDLE = 2000;
var SKIP_COOLDOWN_MAX = 10;
var MAX_SELECTORS_PER_CHECK = 15;

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
	'audio-ads.spotify.com',
	'ads-audio.spotify.com',
	'adserver',
	'app-measurement',
	'analytics',
	'sp-logger',
	'adlog',
	'tracking',
	'metrics',
	'telemetry'
];

var AD_SELECTORS = [
	'[data-testid="ad-type-banner"]',
	'[data-testid="billboard-ad"]',
	'[data-testid="leaderboard-ad"]',
	'[data-testid="sponsorship-ad"]',
	'[data-testid="hpto-ad"]',
	'[data-testid="ad-card"]',
	'[data-testid="ad-slot"]',
	'[data-testid="now-playing-bar-ad"]',
	'[data-testid="inactive-ad"]',
	'[class*="BillboardAd"]',
	'[class*="LeaderboardAd"]',
	'[class*="ad-break"]',
	'[class*="AdBreak"]',
	'[aria-label="Advertisement"]',
	'[aria-label="Publicidad"]',
	'[class*="Commercial"]'
];

var AD_AUDIO_INDICATORS = [
	'[data-testid="ad-type-banner"]',
	'[data-testid="now-playing-bar-ad"]',
	'[class*="Commercial"]',
	'[class*="AdBreak"]',
	'[aria-label*="Advertisement"]',
	'[aria-label*="Publicidad"]',
	'[data-testid="ad-cta"]',
	'[data-testid="ad-message"]'
];

var PROTECTED_SELECTORS = [
	'[data-testid="track-info"]',
	'.track-info',
	'[class*="TrackInfo"]',
	'[data-testid="context-menu"]',
	'[class*="PlaybackControl"]',
	'[class*="PlayerControls"]',
	'[data-testid="control-button"]',
	'[data-testid="volume-bar"]',
	'[class*="VolumeBar"]',
	'[data-testid="extra-controls"]',
	'[class*="ExtraControls"]',
	'.main-nowPlayingBar-nowPlayingBar',
	'.main-nowPlayingBar-NowPlayingBar',
	'[data-testid="now-playing-bar"]',
	'.now-playing-bar',
	'footer[class*="NowPlayingBar"]',
	'footer[class*="nowPlayingBar"]'
];

var NOW_PLAYING_BAR_SELECTORS = [
	'.main-nowPlayingBar-nowPlayingBar',
	'.main-nowPlayingBar-NowPlayingBar',
	'[data-testid="now-playing-bar"]',
	'.now-playing-bar',
	'footer[class*="NowPlayingBar"]',
	'footer[class*="nowPlayingBar"]'
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
		console.log('[Spotilite AdBlock] Blocked ad request:', url);
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
		console.log('[Spotilite AdBlock] Blocked XHR ad request:', url);
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

var debounceTimer = null;
function debounce(func, wait) {
	return function() {
		var context = this;
		var args = arguments;
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(function() {
			func.apply(context, args);
		}, wait);
	};
}

function hideAdElements() {
	try {
		var now = Date.now();
		if (now - window.__spotiliteLastCheck < 500) return;
		window.__spotiliteLastCheck = now;

		var hidden = 0;
		var selectorsToCheck = AD_SELECTORS.slice(0, MAX_SELECTORS_PER_CHECK);
		
		for (var s = 0; s < selectorsToCheck.length; s++) {
			try {
				var sel = selectorsToCheck[s];
				var els = document.querySelectorAll(sel);
				for (var e = 0; e < els.length; e++) {
					var el = els[e];
					if (el.style.display !== 'none') {
						el.style.cssText = 'display:none!important;visibility:hidden!important;opacity:0!important;pointer-events:none!important;height:0!important;width:0!important;overflow:hidden!important;position:absolute!important;left:-9999px!important;top:-9999px!important';
						hidden++;
					}
				}
			} catch(e) {}
		}

		if (hidden > 0 && window.__spotiliteDebug) {
			console.log('[Spotilite AdBlock] Hidden ' + hidden + ' elements');
		}
	} catch(e) {
		console.error('[Spotilite AdBlock] Error hiding elements:', e);
	}
}

function protectPlayerInfo() {
	try {
		// First, ensure all now playing bar containers are visible
		for (var n = 0; n < NOW_PLAYING_BAR_SELECTORS.length; n++) {
			var npb = document.querySelector(NOW_PLAYING_BAR_SELECTORS[n]);
			if (npb) {
				if (npb.style.display === 'none' || npb.style.visibility === 'hidden' || npb.style.opacity === '0') {
					npb.style.cssText = 'display:flex!important;visibility:visible!important;opacity:1!important;pointer-events:auto!important;position:relative!important;height:auto!important;min-height:90px!important;z-index:1000!important';
				}
				// Also protect ALL children of the now playing bar
				var allChildren = npb.querySelectorAll('*');
				for (var c = 0; c < allChildren.length; c++) {
					var child = allChildren[c];
					if (child.style.display === 'none' || child.style.visibility === 'hidden' || child.style.opacity === '0') {
						child.style.cssText = 'display:inherit!important;visibility:visible!important;opacity:1!important;pointer-events:auto!important;position:static!important';
					}
				}
			}
		}
		
		// Then protect specific player elements
		for (var i = 0; i < PROTECTED_SELECTORS.length; i++) {
			var els = document.querySelectorAll(PROTECTED_SELECTORS[i]);
			for (var j = 0; j < els.length; j++) {
				var el = els[j];
				if (el.style.display === 'none' || el.style.visibility === 'hidden' || el.style.opacity === '0') {
					el.style.cssText = 'display:inherit!important;visibility:visible!important;opacity:1!important;pointer-events:auto!important;position:static!important';
				}
			}
		}
	} catch(e) {
		console.error('[Spotilite AdBlock] Error protecting player info:', e);
	}
}

function checkIfAdPlaying() {
	try {
		var nowPlayingBar = document.querySelector('.main-nowPlayingBar-nowPlayingBar, [data-testid="now-playing-bar"], .now-playing-bar');
		if (!nowPlayingBar) return false;

		for (var i = 0; i < AD_AUDIO_INDICATORS.length; i++) {
			var adEl = nowPlayingBar.querySelector(AD_AUDIO_INDICATORS[i]);
			if (adEl && adEl.offsetParent !== null) {
				return true;
			}
		}

		var skipBtn = document.querySelector('[data-testid="skip-button"], [data-testid="control-button-skip-forward"]');
		if (skipBtn) {
			var isVisible = skipBtn.offsetParent !== null && !skipBtn.disabled;
			if (isVisible) {
				var ariaLabel = skipBtn.getAttribute('aria-label') || '';
				if (ariaLabel.toLowerCase().includes('skip') || ariaLabel.toLowerCase().includes('saltar')) {
					var trackInfo = nowPlayingBar.querySelector('[data-testid="track-info"], .track-info');
					if (!trackInfo || !trackInfo.offsetParent) {
						return true;
					}
				}
			}
		}

		var adBadge = nowPlayingBar.querySelector('[data-testid="now-playing-bar-ad"]');
		if (adBadge && adBadge.offsetParent !== null) return true;

		var entityBadges = nowPlayingBar.querySelectorAll('[class*="badge"], [class*="Badge"]');
		for (var j = 0; j < entityBadges.length; j++) {
			var text = entityBadges[j].textContent || '';
			if (text.toLowerCase().includes('ad') || text.toLowerCase().includes('sponsor') || text.toLowerCase().includes('publicidad') || text.toLowerCase().includes('anuncio')) {
				return true;
			}
		}
	} catch(e) {
		console.error('[Spotilite AdBlock] Error checking ad:', e);
	}
	return false;
}

function trySkipAd() {
	if (!window.__spotiliteIsAdPlaying) {
		var isAd = checkIfAdPlaying();
		if (!isAd) return false;
		window.__spotiliteAdConfirmCount++;
		if (window.__spotiliteAdConfirmCount < AD_CONFIRM_THRESHOLD) {
			console.log('[Spotilite AdBlock] Ad detected (' + window.__spotiliteAdConfirmCount + '/' + AD_CONFIRM_THRESHOLD + ')');
			return false;
		}
		console.log('[Spotilite AdBlock] Ad confirmed, attempting skip');
		window.__spotiliteIsAdPlaying = true;
	}

	try {
		var skipBtn = document.querySelector('[data-testid="skip-button"]');
		if (skipBtn && skipBtn.offsetParent !== null && !skipBtn.disabled) {
			skipBtn.click();
			console.log('[Spotilite AdBlock] Skipped ad via skip button');
			resetAdState();
			return true;
		}

		var nextBtn = document.querySelector('[data-testid="control-button-skip-forward"]');
		if (nextBtn && nextBtn.offsetParent !== null && !nextBtn.disabled) {
			nextBtn.click();
			console.log('[Spotilite AdBlock] Skipped ad via next button');
			resetAdState();
			return true;
		}

		var player = document.querySelector('audio[data-testid="audio-player"], .Html5AudioPlayer audio');
		if (player) {
			var currentTime = player.currentTime;
			var duration = player.duration;
			if (duration && duration > 0 && duration < 30 && currentTime < 5) {
				player.currentTime = duration - 0.1;
				console.log('[Spotilite AdBlock] Fast-forwarded short ad');
				resetAdState();
				return true;
			}
		}
	} catch(e) {
		console.error('[Spotilite AdBlock] Error skipping ad:', e);
	}
	return false;
}

function resetAdState() {
	window.__spotiliteAdConfirmCount = 0;
	window.__spotiliteIsAdPlaying = false;
}

var visibilityObserver = null;
function startIntersectionObserver() {
	if (!('IntersectionObserver' in window)) return;

	var adElements = document.querySelectorAll(AD_SELECTORS.join(','));
	adElements.forEach(function(el) {
		visibilityObserver.observe(el);
	});
}

if ('IntersectionObserver' in window) {
	visibilityObserver = new IntersectionObserver(function(entries) {
		entries.forEach(function(entry) {
			if (entry.isIntersecting && entry.target) {
				var el = entry.target;
				for (var s = 0; s < AD_SELECTORS.length; s++) {
					if (el.matches(AD_SELECTORS[s])) {
						el.style.cssText = 'display:none!important;visibility:hidden!important;opacity:0!important;pointer-events:none!important;height:0!important;width:0!important;overflow:hidden!important;position:absolute!important;left:-9999px!important;top:-9999px!important';
						break;
					}
				}
			}
		});
	}, { threshold: 0 });
}

var mutationObserver = null;
var domChangeBuffer = [];
var domChangeTimer = null;

function processDomChanges() {
	if (domChangeBuffer.length === 0) return;
	
	var changes = domChangeBuffer.splice(0);
	var shouldHide = false;
	var shouldCheckSkip = false;

	for (var m = 0; m < changes.length; m++) {
		var mutation = changes[m];
		if (mutation.type === 'childList') {
			var addedNodes = mutation.addedNodes;
			for (var n = 0; n < addedNodes.length; n++) {
				var node = addedNodes[n];
				if (node.nodeType === Node.ELEMENT_NODE) {
					for (var s = 0; s < AD_SELECTORS.length; s++) {
						try {
							if (node.matches && node.matches(AD_SELECTORS[s])) {
								shouldHide = true;
								break;
							}
						} catch(e) {}
					}
				}
			}
		}

		if (mutation.type === 'attributes' && mutation.target.nodeType === Node.ELEMENT_NODE) {
			var attr = mutation.attributeName;
			if (attr === 'class' || attr === 'data-testid' || attr === 'aria-label') {
				var val = mutation.target.getAttribute(attr) || '';
				if (val.toLowerCase().includes('ad') || val.toLowerCase().includes('sponsor') || val.toLowerCase().includes('publicidad')) {
					shouldHide = true;
					shouldCheckSkip = true;
				}
			}
		}
	}

	if (shouldHide) {
		hideAdElements();
		protectPlayerInfo();
	}
	if (shouldCheckSkip) {
		trySkipAd();
	}
}

function startObserving() {
	var target = document.body || document.documentElement;
	if (target) {
		mutationObserver = new MutationObserver(function(mutations) {
			domChangeBuffer = domChangeBuffer.concat(mutations);
			if (domChangeBuffer.length > 50) domChangeBuffer = domChangeBuffer.slice(-50);
			
			clearTimeout(domChangeTimer);
			domChangeTimer = setTimeout(processDomChanges, 100);
		});
		
		mutationObserver.observe(target, {
			childList: true,
			subtree: true,
			attributes: true,
			attributeFilter: ['class', 'data-testid', 'aria-label']
		});
	}
}

var checkInterval = null;
function startPeriodicCheck() {
	var checkCount = 0;
	var skipCooldown = 0;
	var currentInterval = CHECK_INTERVAL_BASE;

	function checkLoop() {
		checkCount++;
		
		var wasIdle = window.__spotiliteIdleTime > 0;
		window.__spotiliteIdleTime++;
		
		if (checkCount % 2 === 0) {
			hideAdElements();
			protectPlayerInfo();
		}

		if (skipCooldown > 0) {
			skipCooldown--;
		}

		if (checkCount % 5 === 0 && skipCooldown === 0) {
			var skipped = trySkipAd();
			if (skipped) {
				skipCooldown = SKIP_COOLDOWN_MAX;
				window.__spotiliteIdleTime = 0;
				currentInterval = CHECK_INTERVAL_BASE;
			}
		}

		if (checkCount % 20 === 0 && !window.__spotiliteIsAdPlaying) {
			window.__spotiliteAdConfirmCount = 0;
		}

		if (window.__spotiliteIdleTime > 10 && currentInterval === CHECK_INTERVAL_BASE) {
			currentInterval = CHECK_INTERVAL_IDLE;
		} else if (window.__spotiliteIdleTime <= 10 && currentInterval === CHECK_INTERVAL_IDLE) {
			currentInterval = CHECK_INTERVAL_BASE;
		}

		if (checkCount > 500) {
			checkCount = 0;
			window.__spotiliteIdleTime = 0;
		}

		checkInterval = setTimeout(checkLoop, currentInterval);
	}

	checkLoop();
}

function cleanup() {
	if (checkInterval) clearTimeout(checkInterval);
	if (mutationObserver) mutationObserver.disconnect();
	if (visibilityObserver) visibilityObserver.disconnect();
	if (domChangeTimer) clearTimeout(domChangeTimer);
}

hideAdElements();
protectPlayerInfo();
startObserving();
startPeriodicCheck();
startIntersectionObserver();

window.__hideAdElements = hideAdElements;
window.__trySkipAd = trySkipAd;
window.__checkIfAdPlaying = checkIfAdPlaying;
window.__protectPlayerInfo = protectPlayerInfo;
window.__cleanup = cleanup;

console.log('[Spotilite AdBlock] Optimized - IntersectionObserver + Debounce + Lazy Loading');
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
		`[data-testid="now-playing-bar-ad"]`,
		`[class*="Commercial"]`,
	}
}