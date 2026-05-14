package modules

const adblockCSS = `
/* Ad containers and banners */
[data-testid="ad-type-banner"],
[data-testid="billboard-ad"],
[data-testid="leaderboard-ad"],
[data-testid="sponsorship-ad"],
[data-testid="hpto-ad"],
[data-testid="ad-card"],
[data-testid="ad-slot"],
[data-testid="now-playing-bar-ad"],
[data-testid="commercial-break"],
[data-testid="ad-overlay"],
[data-testid="inactive-ad"],
[data-testid="ad-cta"],
[data-testid="ad-message"],
[data-testid="sponsored-card"],
[data-testid="ad-break"],
[data-testid="ad-format-video"],
[data-testid="ad-format-audio"],
[data-testid="ad-format-display"],
[data-testid="ad-skip-button-container"],
[data-testid="ad-countdown"],
[data-testid="ad-feedback"],
[data-testid="sponsored-playlist"],
[data-testid="sponsored-content"],
[data-testid="promoted-content"],
[data-testid="ad-label"],
[data-testid="ad-badge"],
[aria-label="Advertisement"],
[aria-label="Publicidad"],
[aria-label="Anuncio"],
[aria-label="Sponsored"],
[aria-label="Patrocinado"],
/* Ad iframes */
iframe[src*="doubleclick.net"],
iframe[src*="moatads"],
iframe[src*="googlesyndication"],
iframe[src*="googleadservices"],
iframe[src*="adservice"],
iframe[src*="ads-"],
/* Legacy class-based selectors */
.WCbmOh4S3HVpA8RhH5Nj,
.Vs2HPUVcMf1MUfOb8KqE,
.main-leaderboardComponent-container,
.main-billboardComponent-container,
.ad-container,
.advertisement,
.sponsored-content,
.promoted-content,
/* Podcast ads */
[data-testid="podcast-ad"],
[data-testid="podcast-sponsor"],
.podcast-ad-marker,
/* Video ads */
[data-testid="video-ad-ui"],
[data-testid="video-ad-container"],
/* Premium upsell (not ads but annoying) */
[data-testid="premium-upsell"],
[data-testid="upgrade-button-container"] {
display: none !important;
visibility: hidden !important;
opacity: 0 !important;
height: 0 !important;
width: 0 !important;
pointer-events: none !important;
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
window.__spotiliteAdSkipAttempts = 0;
window.__spotiliteMutedForAd = false;
window.__spotiliteOriginalVolume = 1;
window.__spotiliteAdCache = {};
window.__spotiliteAdCacheExpiry = 5000;
window.__spotiliteLastAdUrl = null;
window.__spotiliteConsecutiveAdDetections = 0;
window.__spotiliteAdRecoveryAttempts = 0;
window.__spotiliteMaxRecoveryAttempts = 3;

var LOCAL_API = 'localhost:8765';
var AD_CONFIRM_THRESHOLD = 2;
var CHECK_INTERVAL_BASE = 400;
var CHECK_INTERVAL_IDLE = 1000;
var CHECK_INTERVAL_ACTIVE_AD = 200;
var SKIP_COOLDOWN_MAX = 5;
var MAX_AD_DURATION = 60;
var FAST_FORWARD_STEP = 15;

var AD_URL_PATTERNS = [
'pubads.g.doubleclick.net',
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
'ad-delivery.net',
'betweendigital.com',
'tribalfusion.com',
'exponential.com',
'adswizz.com',
'audio-ads.spotify.com',
'ads-audio.spotify.com',
'ad-logger.spotify.com',
'adeventtracker',
'ad-handler.spotify.com',
'spotify.com/ads/',
'spotify.com/audioad',
'spotify.com/ad-logic',
'spotify.com/gabo-receiver-service',
'audio-sp.spotify.com/ad',
'audio-fa.scdn.co/ad',
'heads-fa.scdn.co/ad',
'partnerakamai.spotify.com/ad',
'megaphone.fm/ad',
'audio-ak-spotify-com.akamaized.net/ad',
'audio4-ak-spotify-com.akamaized.net/ad',
'adclick',
'adserver',
'adtech',
'advertising',
'advert',
'/ad/',
'/ads/',
'_ad_',
'_ads_',
'-ad-',
'-ads-',
'ad.doubleclick',
'ad.spotify',
'ads.spotify',
'adlog',
'admetrics',
'adtracking',
'sponsored-audio',
'promo-audio',
'commercial-audio'
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
'[data-testid="ad-cta"]',
'[data-testid="ad-message"]',
'[data-testid="sponsored-card"]',
'[data-testid="ad-break"]',
'[data-testid="commercial-break"]',
'[data-testid="ad-overlay"]',
'[data-testid="ad-format-video"]',
'[data-testid="ad-format-audio"]',
'[data-testid="ad-format-display"]',
'[data-testid="ad-skip-button-container"]',
'[data-testid="ad-countdown"]',
'[data-testid="ad-feedback"]',
'[data-testid="sponsored-playlist"]',
'[data-testid="sponsored-content"]',
'[data-testid="promoted-content"]',
'[data-testid="ad-label"]',
'[data-testid="ad-badge"]',
'[data-testid="podcast-ad"]',
'[data-testid="podcast-sponsor"]',
'[data-testid="video-ad-ui"]',
'[data-testid="video-ad-container"]',
'[aria-label="Advertisement"]',
'[aria-label="Publicidad"]',
'[aria-label="Anuncio"]',
'[aria-label="Sponsored"]',
'[aria-label="Patrocinado"]',
'.ad-container',
'.advertisement',
'.sponsored-content',
'.promoted-content',
'.podcast-ad-marker'
];

var AD_AUDIO_INDICATORS = [
'[data-testid="ad-type-banner"]',
'[data-testid="now-playing-bar-ad"]',
'[data-testid="ad-cta"]',
'[data-testid="ad-message"]',
'[data-testid="commercial-break"]',
'[data-testid="ad-break"]',
'[data-testid="ad-overlay"]',
'[data-testid="ad-format-audio"]',
'[data-testid="ad-countdown"]',
'[data-testid="ad-label"]',
'[data-testid="podcast-ad"]',
'[aria-label="Advertisement"]',
'[aria-label="Publicidad"]',
'[aria-label="Anuncio"]'
];

var AD_TEXT_PATTERNS = [
'advertisement',
'publicidad',
'anuncio',
'sponsored',
'patrocinado',
'commercial break',
'pausa comercial',
'ad playing',
'reproduciendo anuncio'
];

var PROTECTED_SELECTORS = [
'[data-testid="track-info"]',
'.track-info',
'[class*="TrackInfo"]',
'[data-testid="context-menu"]',
'[class*="PlaybackControl"]',
'[data-testid="control-button"]',
'[class*="NowPlayingBar"]',
'[data-testid="track-title"]',
'[data-testid="track-artist"]',
'[class*="TrackTitle"]',
'[class*="ArtistLink"]',
'[data-testid="player-controls"]',
'[class*="PlayerControls"]',
'[data-testid="playback-bar"]',
'[data-testid="playback-bar-inner"]',
'[data-testid="control-button-skip-forward"]',
'[data-testid="control-button-skip-back"]',
'[data-testid="control-button-play-pause"]',
'[data-testid="now-playing-bar"]',
'[data-testid="cover-art-image"]',
'[data-testid="artist-avatar"]',
'[data-testid="entity-image"]'
];

function isAdUrl(url) {
if (!url) return false;
var u = typeof url === 'string' ? url : (url.url ? url.url : '');
if (!u || u.indexOf(LOCAL_API) !== -1) return false;

var now = Date.now();
if (window.__spotiliteAdCache[u]) {
if (now - window.__spotiliteAdCache[u].time < window.__spotiliteAdCacheExpiry) {
return window.__spotiliteAdCache[u].isAd;
}
}

if (u.indexOf('spotify.com') !== -1 && u.indexOf('/ad') === -1 && u.indexOf('/ads') === -1 && u.indexOf('ad-') === -1 && u.indexOf('adlog') === -1) {
window.__spotiliteAdCache[u] = { isAd: false, time: now };
return false;
}

for (var i = 0; i < AD_URL_PATTERNS.length; i++) {
if (u.indexOf(AD_URL_PATTERNS[i]) !== -1) {
window.__spotiliteAdCache[u] = { isAd: true, time: now };
return true;
}
}

window.__spotiliteAdCache[u] = { isAd: false, time: now };
return false;
}

function createEmptyResponse() {
return new Response('{}', {
status: 200,
headers: { 'Content-Type': 'application/json' }
});
}

function createEmptyAudioResponse() {
return new Response(new ArrayBuffer(0), {
status: 200,
headers: { 'Content-Type': 'audio/mpeg' }
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
if (body && (body.type === 'ad' || body.typ === 'ad' || body.ad || body.isAd || body.ad_break || body.adbreak)) {
return Promise.resolve(createEmptyResponse());
}
} catch(e) {}
}

var originalFetch = window.__origFetch.apply(this, args);

if (originalFetch && originalFetch.then) {
return originalFetch.then(function(response) {
try {
var contentType = response.headers.get('content-type') || '';
var responseUrl = response.url || '';

if (isAdUrl(responseUrl)) {
console.log('[Spotilite AdBlock] Blocked ad response:', responseUrl);
if (contentType.indexOf('audio') !== -1) {
return createEmptyAudioResponse();
}
return createEmptyResponse();
}

if (contentType.indexOf('audio') !== -1 || responseUrl.indexOf('.mp3') !== -1 || responseUrl.indexOf('audio') !== -1) {
for (var i = 0; i < AD_URL_PATTERNS.length; i++) {
if (responseUrl.indexOf(AD_URL_PATTERNS[i]) !== -1) {
console.log('[Spotilite AdBlock] Blocked audio ad:', responseUrl);
return createEmptyAudioResponse();
}
}
}

if (contentType.indexOf('json') !== -1) {
var clone = response.clone();
return clone.json().then(function(data) {
if (data && (data.type === 'ad' || data.ad || data.isAd || data.ad_break || data.adbreak || data.ad_metadata || data.ad_info)) {
console.log('[Spotilite AdBlock] Blocked ad in JSON response');
return createEmptyResponse();
}
if (data && data.content && data.ad) {
return createEmptyResponse();
}
return response;
}).catch(function() {
return response;
});
}
} catch(e) {}
return response;
});
}

return originalFetch;
};

var _origXHROpen = XMLHttpRequest.prototype.open;
var _origXHRSend = XMLHttpRequest.prototype.send;

XMLHttpRequest.prototype.open = function(method, url) {
this.__requestUrl = url;
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
var self = this;
setTimeout(function() {
self.readyState = 4;
self.status = 200;
self.responseText = '{}';
self.onreadystatechange && self.onreadystatechange();
self.onload && self.onload();
}, 0);
return;
}
return _origXHRSend.apply(this, arguments);
};

var _origXHRAddEventListener = XMLHttpRequest.prototype.addEventListener;
XMLHttpRequest.prototype.addEventListener = function(event, callback) {
if (this.__adBlocked && (event === 'load' || event === 'readystatechange')) {
return;
}
return _origXHRAddEventListener.apply(this, arguments);
};

function hideAdElements() {
try {
var now = Date.now();
if (now - window.__spotiliteLastCheck < 300) return;
window.__spotiliteLastCheck = now;

var hidden = 0;
var highPriorityHidden = 0;

var highPrioritySelectors = [
'[data-testid="commercial-break"]',
'[data-testid="ad-overlay"]',
'[data-testid="ad-format-audio"]',
'[data-testid="ad-format-video"]'
];

for (var s = 0; s < highPrioritySelectors.length; s++) {
try {
var els = document.querySelectorAll(highPrioritySelectors[s]);
for (var e = 0; e < els.length; e++) {
var el = els[e];
if (el.style.display !== 'none') {
el.style.cssText = 'display:none!important;visibility:hidden!important;opacity:0!important;height:0!important;width:0!important;pointer-events:none!important';
highPriorityHidden++;
hidden++;
}
}
} catch(e) {}
}

for (var s = 0; s < AD_SELECTORS.length; s++) {
try {
var els = document.querySelectorAll(AD_SELECTORS[s]);
for (var e = 0; e < els.length; e++) {
var el = els[e];
if (el.style.display !== 'none') {
el.style.display = 'none';
hidden++;
}
}
} catch(e) {}
}

if (highPriorityHidden > 0) {
console.log('[Spotilite AdBlock] Blocked ' + highPriorityHidden + ' high-priority ad elements');
window.__spotiliteConsecutiveAdDetections++;
} else if (hidden === 0) {
window.__spotiliteConsecutiveAdDetections = 0;
}

if (hidden > 0 && window.__spotiliteDebug) {
console.log('[Spotilite AdBlock] Hidden ' + hidden + ' ad elements');
}
} catch(e) {
console.error('[Spotilite AdBlock] Error hiding elements:', e);
}
}

function protectPlayerInfo() {
try {
for (var i = 0; i < PROTECTED_SELECTORS.length; i++) {
var els = document.querySelectorAll(PROTECTED_SELECTORS[i]);
for (var j = 0; j < els.length; j++) {
var el = els[j];
if (el.style.display === 'none' || el.style.visibility === 'hidden') {
el.style.cssText = 'display:inherit!important;visibility:visible!important;opacity:1!important;pointer-events:auto!important;position:static!important';
}
}
}
} catch(e) {}
}

function checkIfAdPlaying() {
try {
var nowPlayingBar = document.querySelector('.main-nowPlayingBar-nowPlayingBar, [data-testid="now-playing-bar"], .now-playing-bar');
if (!nowPlayingBar) return false;

for (var i = 0; i < AD_AUDIO_INDICATORS.length; i++) {
var adEl = nowPlayingBar.querySelector(AD_AUDIO_INDICATORS[i]);
if (adEl && adEl.offsetParent !== null) {
console.log('[Spotilite AdBlock] Ad detected via indicator:', AD_AUDIO_INDICATORS[i]);
return true;
}
}

var adElGlobal = document.querySelector('[data-testid="commercial-break"], [data-testid="ad-break"], [data-testid="ad-overlay"], [data-testid="ad-format-audio"]');
if (adElGlobal && adElGlobal.offsetParent !== null) {
console.log('[Spotilite AdBlock] Ad detected via global element');
return true;
}

var textContent = nowPlayingBar.textContent || '';
var lowerText = textContent.toLowerCase();
for (var p = 0; p < AD_TEXT_PATTERNS.length; p++) {
if (lowerText.indexOf(AD_TEXT_PATTERNS[p]) !== -1) {
console.log('[Spotilite AdBlock] Ad detected via text pattern:', AD_TEXT_PATTERNS[p]);
return true;
}
}

var trackInfo = nowPlayingBar.querySelector('[data-testid="track-info"], .track-info, [class*="TrackInfo"]');
var hasTrackInfo = trackInfo && trackInfo.offsetParent !== null;

var titleEl = nowPlayingBar.querySelector('[data-testid="track-title"], .track-title, [class*="TrackTitle"]');
var hasTitle = titleEl && titleEl.textContent && titleEl.textContent.trim().length > 0;

var artistEl = nowPlayingBar.querySelector('[data-testid="context-item-info-artist"], [class*="ArtistLink"], [data-testid="track-artist"]');
var hasArtist = artistEl && artistEl.offsetParent !== null;

if (!hasTrackInfo && !hasTitle && !hasArtist) {
var skipBtn = document.querySelector('[data-testid="skip-button"], [data-testid="control-button-skip-forward"]');
if (skipBtn) {
var isVisible = skipBtn.offsetParent !== null && !skipBtn.disabled;
if (isVisible) {
console.log('[Spotilite AdBlock] Ad detected via missing track info + visible skip button');
return true;
}
}
}

var adBadge = nowPlayingBar.querySelector('[data-testid="now-playing-bar-ad"], [data-testid="ad-badge"], [data-testid="ad-label"]');
if (adBadge && adBadge.offsetParent !== null) {
console.log('[Spotilite AdBlock] Ad detected via ad badge');
return true;
}

var playerControls = document.querySelector('[data-testid="player-controls"], [class*="PlayerControls"]');
if (playerControls) {
var adText = playerControls.textContent || '';
if (adText.toLowerCase().includes('ad ') || adText.toLowerCase().includes('anuncio') || adText.toLowerCase().includes('commercial')) {
console.log('[Spotilite AdBlock] Ad detected via player controls text');
return true;
}
}

var player = getPlayerElement();
if (player && player.src) {
if (isAdUrl(player.src)) {
console.log('[Spotilite AdBlock] Ad detected via audio source URL');
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
if (!isAd) {
if (window.__spotiliteMutedForAd) {
unmuteAudio();
}
window.__spotiliteAdRecoveryAttempts = 0;
return false;
}
window.__spotiliteAdConfirmCount++;
if (window.__spotiliteAdConfirmCount < AD_CONFIRM_THRESHOLD) {
console.log('[Spotilite AdBlock] Ad detected (' + window.__spotiliteAdConfirmCount + '/' + AD_CONFIRM_THRESHOLD + ')');
return false;
}
console.log('[Spotilite AdBlock] Ad confirmed, initiating skip sequence');
window.__spotiliteIsAdPlaying = true;
window.__spotiliteAdSkipAttempts = 0;
window.__spotiliteAdRecoveryAttempts = 0;
}

window.__spotiliteAdSkipAttempts++;

try {
muteAudioForAd();

var skipBtn = document.querySelector('[data-testid="control-button-skip-forward"], [data-testid="skip-button"], button[aria-label*="Skip"], button[aria-label*="Saltar"], button[aria-label*="Next"]');
if (skipBtn && skipBtn.offsetParent !== null && !skipBtn.disabled) {
if (window.__spotiliteAdSkipAttempts <= 5) {
skipBtn.click();
console.log('[Spotilite AdBlock] Clicked skip button (attempt ' + window.__spotiliteAdSkipAttempts + ')');
}
}

var player = getPlayerElement();
if (player) {
var currentTime = player.currentTime;
var duration = player.duration;
if (duration && duration > 0 && duration < MAX_AD_DURATION) {
if (currentTime < duration - 1) {
var newTime = Math.min(duration - 0.3, currentTime + FAST_FORWARD_STEP);
player.currentTime = newTime;
console.log('[Spotilite AdBlock] Fast-forwarded ad: ' + currentTime.toFixed(1) + 's -> ' + newTime.toFixed(1) + 's / ' + duration.toFixed(1) + 's');
}
}

if (isAdUrl(player.src)) {
console.log('[Spotilite AdBlock] Detected ad audio source, attempting to clear');
player.pause();
player.src = '';
player.load();
}
}

if (window.__spotiliteAdSkipAttempts > 30) {
console.log('[Spotilite AdBlock] Ad skip timeout, attempting recovery');
attemptAdRecovery();
return true;
}

if (window.__spotiliteAdSkipAttempts > 10 && window.__spotiliteAdSkipAttempts % 5 === 0) {
if (!checkIfAdPlaying()) {
console.log('[Spotilite AdBlock] Ad appears to have ended');
resetAdState();
return true;
}
}
} catch(e) {
console.error('[Spotilite AdBlock] Error handling ad:', e);
}
return false;
}

function attemptAdRecovery() {
if (window.__spotiliteAdRecoveryAttempts >= window.__spotiliteMaxRecoveryAttempts) {
console.log('[Spotilite AdBlock] Max recovery attempts reached, forcing reset');
resetAdState();
return;
}

window.__spotiliteAdRecoveryAttempts++;
console.log('[Spotilite AdBlock] Recovery attempt ' + window.__spotiliteAdRecoveryAttempts);

try {
var player = getPlayerElement();
if (player) {
player.pause();
setTimeout(function() {
var nextBtn = document.querySelector('[data-testid="control-button-skip-forward"], button[aria-label*="Next"]');
if (nextBtn && !nextBtn.disabled) {
nextBtn.click();
console.log('[Spotilite AdBlock] Clicked next track button for recovery');
}
resetAdState();
}, 500);
}
} catch(e) {
console.error('[Spotilite AdBlock] Recovery error:', e);
resetAdState();
}
}

function getPlayerElement() {
var player = document.querySelector('audio[data-testid="audio-player"]');
if (player) return player;

player = document.querySelector('.Html5AudioPlayer audio');
if (player) return player;

player = document.querySelector('audio[src*="spotify"]');
if (player) return player;

player = document.querySelector('audio[src*="scdn.co"]');
if (player) return player;

player = document.querySelector('audio[src*="akamaized.net"]');
if (player) return player;

var audios = document.querySelectorAll('audio');
for (var i = 0; i < audios.length; i++) {
if (audios[i].src && (audios[i].src.indexOf('spotify') !== -1 || audios[i].src.indexOf('scdn.co') !== -1 || audios[i].src.indexOf('akamaized') !== -1)) {
return audios[i];
}
}

if (audios.length > 0) return audios[0];

return null;
}

function muteAudioForAd() {
if (window.__spotiliteMutedForAd) return;

var player = getPlayerElement();
if (player) {
window.__spotiliteOriginalVolume = player.volume;
player.volume = 0;
player.muted = true;
window.__spotiliteMutedForAd = true;
console.log('[Spotilite AdBlock] Audio muted for ad');
}
}

function unmuteAudio() {
if (!window.__spotiliteMutedForAd) return;

var player = getPlayerElement();
if (player) {
player.volume = window.__spotiliteOriginalVolume;
player.muted = false;
}

var allAudios = document.querySelectorAll('audio');
for (var i = 0; i < allAudios.length; i++) {
allAudios[i].muted = false;
}

window.__spotiliteMutedForAd = false;
console.log('[Spotilite AdBlock] Audio unmuted');
}

function resetAdState() {
window.__spotiliteAdConfirmCount = 0;
window.__spotiliteIsAdPlaying = false;
window.__spotiliteAdSkipAttempts = 0;
window.__spotiliteConsecutiveAdDetections = 0;
window.__spotiliteAdRecoveryAttempts = 0;
if (window.__spotiliteMutedForAd) {
unmuteAudio();
}
console.log('[Spotilite AdBlock] Ad state reset');
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
if (attr === 'data-testid') {
var val = mutation.target.getAttribute('data-testid') || '';
if (val.indexOf('ad-') === 0 || val === 'commercial-break' || val === 'ad-overlay' || val === 'sponsored-card') {
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
if (domChangeBuffer.length > 100) domChangeBuffer = domChangeBuffer.slice(-100);

clearTimeout(domChangeTimer);
domChangeTimer = setTimeout(processDomChanges, 50);
});

mutationObserver.observe(target, {
childList: true,
subtree: true,
attributes: true,
attributeFilter: ['data-testid', 'aria-label']
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

window.__spotiliteIdleTime++;

hideAdElements();
protectPlayerInfo();

if (skipCooldown > 0) {
skipCooldown--;
}

var isActiveAd = window.__spotiliteIsAdPlaying || window.__spotiliteConsecutiveAdDetections > 2;
if (isActiveAd) {
currentInterval = CHECK_INTERVAL_ACTIVE_AD;
window.__spotiliteIdleTime = 0;
}

if (checkCount % 2 === 0 || isActiveAd) {
if (skipCooldown === 0) {
var skipped = trySkipAd();
if (skipped) {
skipCooldown = SKIP_COOLDOWN_MAX;
window.__spotiliteIdleTime = 0;
currentInterval = CHECK_INTERVAL_BASE;
}
}
}

if (checkCount % 15 === 0 && !window.__spotiliteIsAdPlaying) {
window.__spotiliteAdConfirmCount = Math.max(0, window.__spotiliteAdConfirmCount - 1);
}

if (checkCount % 50 === 0) {
var cacheSize = Object.keys(window.__spotiliteAdCache).length;
if (cacheSize > 100) {
console.log('[Spotilite AdBlock] Clearing ad URL cache (' + cacheSize + ' entries)');
window.__spotiliteAdCache = {};
}
}

if (window.__spotiliteIdleTime > 15 && currentInterval === CHECK_INTERVAL_BASE) {
currentInterval = CHECK_INTERVAL_IDLE;
} else if (window.__spotiliteIdleTime <= 15 && currentInterval === CHECK_INTERVAL_IDLE && !isActiveAd) {
currentInterval = CHECK_INTERVAL_BASE;
}

if (checkCount > 10000) {
checkCount = 0;
window.__spotiliteIdleTime = Math.min(window.__spotiliteIdleTime, 100);
}

checkInterval = setTimeout(checkLoop, currentInterval);
}

checkLoop();
}

function cleanup() {
if (checkInterval) clearTimeout(checkInterval);
if (mutationObserver) mutationObserver.disconnect();
if (domChangeTimer) clearTimeout(domChangeTimer);
}

hideAdElements();
protectPlayerInfo();
startObserving();
startPeriodicCheck();

window.__hideAdElements = hideAdElements;
window.__trySkipAd = trySkipAd;
window.__checkIfAdPlaying = checkIfAdPlaying;
window.__protectPlayerInfo = protectPlayerInfo;
window.__cleanup = cleanup;
window.__muteAudioForAd = muteAudioForAd;
window.__unmuteAudio = unmuteAudio;
window.__resetAdState = resetAdState;
window.__attemptAdRecovery = attemptAdRecovery;
window.__isAdUrl = isAdUrl;

console.log('[Spotilite AdBlock] Enhanced mode initialized - Network blocking + Smart ad detection + Auto-skip + Recovery');
console.log('[Spotilite AdBlock] Features: URL caching, priority blocking, text pattern detection, auto-recovery');
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

func (m *AdBlockModule) CSS() string { return adblockCSS }
func (m *AdBlockModule) JS() string  { return adblockJS }
func (m *AdBlockModule) Selectors() []string {
	return []string{
		`[data-testid="ad-type-banner"]`,
		`[data-testid="billboard-ad"]`,
		`[data-testid="leaderboard-ad"]`,
		`[data-testid="sponsorship-ad"]`,
		`[data-testid="hpto-ad"]`,
		`[data-testid="commercial-break"]`,
		`[data-testid="ad-overlay"]`,
		`[data-testid="now-playing-bar-ad"]`,
		`[data-testid="ad-card"]`,
		`[data-testid="ad-slot"]`,
		`[data-testid="ad-format-video"]`,
		`[data-testid="ad-format-audio"]`,
		`[data-testid="ad-format-display"]`,
		`[data-testid="ad-countdown"]`,
		`[data-testid="sponsored-content"]`,
		`[data-testid="promoted-content"]`,
		`[data-testid="podcast-ad"]`,
		`[data-testid="video-ad-container"]`,
		`[aria-label="Advertisement"]`,
		`[aria-label="Publicidad"]`,
		`[aria-label="Anuncio"]`,
		`[aria-label="Sponsored"]`,
		`.ad-container`,
		`.advertisement`,
	}
}
