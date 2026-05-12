package modules

const adblockCSS = `
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
[aria-label="Advertisement"],
[aria-label="Publicidad"],
iframe[src*="doubleclick.net"],
iframe[src*="moatads"],
.WCbmOh4S3HVpA8RhH5Nj,
.Vs2HPUVcMf1MUfOb8KqE,
.main-leaderboardComponent-container,
.main-billboardComponent-container {
	display: none !important;
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

var LOCAL_API = 'localhost:8765';
var AD_CONFIRM_THRESHOLD = 1;
var CHECK_INTERVAL_BASE = 500;
var CHECK_INTERVAL_IDLE = 1000;
var SKIP_COOLDOWN_MAX = 5;
var MAX_SELECTORS_PER_CHECK = 20;

var AD_URL_PATTERNS = [
	'ad-handler.spotify.com',
	'ad-return-url',
	'gaia.spotify.com/ad',
	'spclient.wg.spotify.com/ad',
	'spclient.wg.spotify.com/ad-logic',
	'spclient.wg.spotify.com/ads/',
	'spclient.wg.spotify.com/ad/',
	'spclient.wg.spotify.com/gabo-receiver',
	'partnerakamai.spotify.com/ad',
	'/ad/',
	'/ads/',
	'ads.php',
	'spotify.com/ad-logic',
	'spotify.com/ads/',
	'spotify.com/pair/',
	'spotify.com/gabo-receiver-service',
	'audio-ads.spotify.com',
	'ads-audio.spotify.com',
	'audio-sp.spotify.com/ad',
	'audio-fa.scdn.co/ad',
	'heads-fa.scdn.co/ad',
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
	'adserver',
	'app-measurement',
	'analytics',
	'sp-logger',
	'adlog',
	'tracking',
	'metrics',
	'telemetry',
	'ad-delivery.net',
	'betweendigital.com',
	'tribalfusion.com',
	'exponential.com',
	'spotify.ads',
	'/ad-break',
	'adbreak',
	'ad-content',
	'/commercial/',
	'adeventtracker',
	'adswizz.com',
	'megaphone.fm/ad',
	'spotify.com/audioad',
	'/audioad/',
	'ad-logger.spotify.com',
	'spclient.wg.spotify.com/comscore/',
	'spclient.wg.spotify.com/event/',
	'audio-ak-spotify-com.akamaized.net/ad',
	'audio4-ak-spotify-com.akamaized.net/ad'
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
	'[class*="BillboardAd"]',
	'[class*="LeaderboardAd"]',
	'[class*="ad-break"]',
	'[class*="AdBreak"]',
	'[class*="adBreak"]',
	'[class*="AdBreak"]',
	'[class*="commercial-break"]',
	'[class*="CommercialBreak"]',
	'[aria-label="Advertisement"]',
	'[aria-label="Publicidad"]',
	'[aria-label*="Advertisement"]',
	'[aria-label*="Publicidad"]',
	'[aria-label*="Sponsored"]',
	'[aria-label*="Patrocinado"]',
	'[class*="Commercial"]',
	'[class*="commercial"]',
	'[class*="Sponsored"]',
	'[class*="sponsored"]',
	'[class*="AdOverlay"]',
	'[class*="adOverlay"]',
	'[class*="AdMessage"]',
	'[class*="adMessage"]',
	'[class*="AdBanner"]',
	'[class*="adBanner"]',
	'[class*="AdContainer"]',
	'[class*="adContainer"]',
	'[class*="NowPlayingAd"]',
	'[class*="nowPlayingAd"]'
];

var AD_AUDIO_INDICATORS = [
	'[data-testid="ad-type-banner"]',
	'[data-testid="now-playing-bar-ad"]',
	'[data-testid="ad-cta"]',
	'[data-testid="ad-message"]',
	'[data-testid="commercial-break"]',
	'[data-testid="ad-break"]',
	'[data-testid="ad-overlay"]',
	'[class*="Commercial"]',
	'[class*="commercial"]',
	'[class*="AdBreak"]',
	'[class*="adBreak"]',
	'[class*="CommercialBreak"]',
	'[class*="commercialBreak"]',
	'[aria-label*="Advertisement"]',
	'[aria-label*="Publicidad"]',
	'[aria-label*="Sponsored"]',
	'[aria-label*="Patrocinado"]',
	'[data-testid="ad-cta"]',
	'[data-testid="ad-message"]',
	'[class*="AdMessage"]',
	'[class*="adMessage"]',
	'[class*="NowPlayingAd"]',
	'[class*="nowPlayingAd"]',
	'[class*="AdBanner"]',
	'[class*="adBanner"]'
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
	'[data-testid="now-playing-bar"]'
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
				var clone = response.clone();
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
					return clone.json().then(function(data) {
						if (data && (data.type === 'ad' || data.ad || data.isAd || data.ad_break || data.adbreak || data.ad_metadata || data.ad_info)) {
							console.log('[Spotilite AdBlock] Blocked ad in JSON response');
							return createEmptyResponse();
						}
						if (data && data.content && data.content.ad) {
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
				return true;
			}
		}

		var adElGlobal = document.querySelector('[data-testid="commercial-break"], [data-testid="ad-break"], [data-testid="ad-overlay"], [class*="CommercialBreak"], [class*="AdBreak"], [class*="AdOverlay"]');
		if (adElGlobal && adElGlobal.offsetParent !== null) {
			return true;
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
					return true;
				}
			}
		}

		var adBadge = nowPlayingBar.querySelector('[data-testid="now-playing-bar-ad"], [class*="AdBadge"], [class*="adBadge"], [class*="AdLabel"], [class*="adLabel"]');
		if (adBadge && adBadge.offsetParent !== null) return true;

		var entityBadges = nowPlayingBar.querySelectorAll('[class*="badge"], [class*="Badge"], [class*="label"], [class*="Label"]');
		for (var j = 0; j < entityBadges.length; j++) {
			var text = entityBadges[j].textContent || '';
			if (text.toLowerCase().includes('ad') || text.toLowerCase().includes('sponsor') || text.toLowerCase().includes('publicidad') || text.toLowerCase().includes('anuncio') || text.toLowerCase().includes('commercial') || text.toLowerCase().includes('patrocin')) {
				return true;
			}
		}

		var adLinks = nowPlayingBar.querySelectorAll('a');
		for (var k = 0; k < adLinks.length; k++) {
			var linkText = adLinks[k].textContent || '';
			var href = adLinks[k].href || '';
			if (linkText.toLowerCase().includes('learn more') || linkText.toLowerCase().includes('saber más') || linkText.toLowerCase().includes('conocer más') || href.indexOf('ad') !== -1 || href.indexOf('sponsor') !== -1) {
				return true;
			}
		}

		var playerControls = document.querySelector('[data-testid="player-controls"], [class*="PlayerControls"]');
		if (playerControls) {
			var adText = playerControls.textContent || '';
			if (adText.toLowerCase().includes('ad ') || adText.toLowerCase().includes('anuncio') || adText.toLowerCase().includes('commercial')) {
				return true;
			}
		}

		var playbackBar = document.querySelector('[data-testid="playback-bar"], [data-testid="playback-bar-inner"]');
		if (playbackBar) {
			var barText = playbackBar.textContent || '';
			if (barText.toLowerCase().includes('ad ') || barText.toLowerCase().includes('anuncio')) {
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
			return false;
		}
		window.__spotiliteAdConfirmCount++;
		if (window.__spotiliteAdConfirmCount < AD_CONFIRM_THRESHOLD) {
			console.log('[Spotilite AdBlock] Ad detected (' + window.__spotiliteAdConfirmCount + '/' + AD_CONFIRM_THRESHOLD + ')');
			return false;
		}
		console.log('[Spotilite AdBlock] Ad confirmed, muting audio');
		window.__spotiliteIsAdPlaying = true;
		window.__spotiliteAdSkipAttempts = 0;
	}

	window.__spotiliteAdSkipAttempts++;

	try {
		muteAudioForAd();

		var skipBtn = document.querySelector('[data-testid="control-button-skip-forward"], [data-testid="skip-button"], button[aria-label*="Skip"], button[aria-label*="Saltar"]');
		if (skipBtn && skipBtn.offsetParent !== null && !skipBtn.disabled && window.__spotiliteAdSkipAttempts <= 3) {
			skipBtn.click();
			console.log('[Spotilite AdBlock] Attempted skip via skip button');
		}

		var player = getPlayerElement();
		if (player) {
			var currentTime = player.currentTime;
			var duration = player.duration;
			if (duration && duration > 0 && duration < 60) {
				if (currentTime < 5) {
					player.currentTime = Math.min(duration - 0.5, currentTime + 10);
					console.log('[Spotilite AdBlock] Fast-forwarded ad segment');
				}
			}
		}

		var playbackBar = document.querySelector('[data-testid="playback-bar"], [data-testid="playback-bar-inner"]');
		if (playbackBar && window.__spotiliteAdSkipAttempts <= 2) {
			var barRect = playbackBar.getBoundingClientRect();
			if (barRect.width > 0) {
				var clickX = barRect.left + (barRect.width * 0.95);
				var clickY = barRect.top + (barRect.height / 2);
				var clickEvent = new MouseEvent('click', {
					bubbles: true,
					cancelable: true,
					view: window,
					clientX: clickX,
					clientY: clickY
				});
				playbackBar.dispatchEvent(clickEvent);
				console.log('[Spotilite AdBlock] Attempted seek via playback bar');
			}
		}

		if (window.__spotiliteAdSkipAttempts > 20) {
			console.log('[Spotilite AdBlock] Ad timeout, resetting');
			resetAdState();
			return true;
		}
	} catch(e) {
		console.error('[Spotilite AdBlock] Error handling ad:', e);
	}
	return false;
}
			return false;
		}
		window.__spotiliteAdConfirmCount++;
		if (window.__spotiliteAdConfirmCount < AD_CONFIRM_THRESHOLD) {
			console.log('[Spotilite AdBlock] Ad detected (' + window.__spotiliteAdConfirmCount + '/' + AD_CONFIRM_THRESHOLD + ')');
			return false;
		}
		console.log('[Spotilite AdBlock] Ad confirmed, muting audio');
		window.__spotiliteIsAdPlaying = true;
		window.__spotiliteAdSkipAttempts = 0;
	}

	window.__spotiliteAdSkipAttempts++;

	try {
		muteAudioForAd();

		var skipBtn = document.querySelector('[data-testid="skip-button"]');
		if (skipBtn && skipBtn.offsetParent !== null && !skipBtn.disabled && window.__spotiliteAdSkipAttempts <= 2) {
			skipBtn.click();
			console.log('[Spotilite AdBlock] Attempted skip via skip button');
		}

		var player = getPlayerElement();
		if (player) {
			var currentTime = player.currentTime;
			var duration = player.duration;
			if (duration && duration > 0 && duration < 45) {
				if (currentTime < 3) {
					player.currentTime = Math.min(duration - 0.5, currentTime + 5);
					console.log('[Spotilite AdBlock] Fast-forwarded ad segment');
				}
			}
		}

		if (window.__spotiliteAdSkipAttempts > 15) {
			console.log('[Spotilite AdBlock] Ad timeout, resetting');
			resetAdState();
			return true;
		}
	} catch(e) {
		console.error('[Spotilite AdBlock] Error handling ad:', e);
	}
	return false;
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

function muteAllAudio() {
	var muted = [];
	var audios = document.querySelectorAll('audio');
	for (var i = 0; i < audios.length; i++) {
		var audio = audios[i];
		if (audio.src && !audio.muted) {
			muted.push({el: audio, volume: audio.volume, muted: audio.muted});
			audio.volume = 0;
			audio.muted = true;
		}
	}
	return muted;
}

function unmuteAllAudio(mutedList) {
	if (!mutedList) return;
	for (var i = 0; i < mutedList.length; i++) {
		var item = mutedList[i];
		if (item.el) {
			item.el.volume = item.volume;
			item.el.muted = item.muted;
		}
	}
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
	
	var allMuted = muteAllAudio();
	if (allMuted.length > 0) {
		window.__spotiliteMutedForAd = true;
		console.log('[Spotilite AdBlock] All audio muted for ad (' + allMuted.length + ' sources)');
	}
	
	var rootEl = document.querySelector('.Root__main-view, .main-view-container');
	if (rootEl) {
		rootEl.style.filter = 'grayscale(0.5) brightness(0.7)';
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
	
	var rootEl = document.querySelector('.Root__main-view, .main-view-container');
	if (rootEl) {
		rootEl.style.filter = 'none';
	}
	
	window.__spotiliteMutedForAd = false;
	console.log('[Spotilite AdBlock] Audio unmuted');
}

function resetAdState() {
	window.__spotiliteAdConfirmCount = 0;
	window.__spotiliteIsAdPlaying = false;
	window.__spotiliteAdSkipAttempts = 0;
	if (window.__spotiliteMutedForAd) {
		unmuteAudio();
	}
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
			if (domChangeBuffer.length > 100) domChangeBuffer = domChangeBuffer.slice(-100);
			
			clearTimeout(domChangeTimer);
			domChangeTimer = setTimeout(processDomChanges, 50);
		});
		
		mutationObserver.observe(target, {
			childList: true,
			subtree: true,
			attributes: true,
			attributeFilter: ['class', 'data-testid', 'aria-label', 'style']
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
		
		hideAdElements();
		protectPlayerInfo();

		if (skipCooldown > 0) {
			skipCooldown--;
		}

		if (checkCount % 2 === 0 || window.__spotiliteIsAdPlaying) {
			if (skipCooldown === 0) {
				var skipped = trySkipAd();
				if (skipped) {
					skipCooldown = SKIP_COOLDOWN_MAX;
					window.__spotiliteIdleTime = 0;
					currentInterval = CHECK_INTERVAL_BASE;
				}
			}
		}

		if (checkCount % 10 === 0 && !window.__spotiliteIsAdPlaying) {
			window.__spotiliteAdConfirmCount = 0;
		}

		if (window.__spotiliteIdleTime > 10 && currentInterval === CHECK_INTERVAL_BASE) {
			currentInterval = CHECK_INTERVAL_IDLE;
		} else if (window.__spotiliteIdleTime <= 10 && currentInterval === CHECK_INTERVAL_IDLE) {
			currentInterval = CHECK_INTERVAL_BASE;
		}

		if (checkCount > 1000) {
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
window.__muteAudioForAd = muteAudioForAd;
window.__unmuteAudio = unmuteAudio;

console.log('[Spotilite AdBlock] Aggressive mode - Mute on ad + Fast detection + Network block');
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
		`[data-testid="commercial-break"]`,
		`[data-testid="ad-overlay"]`,
		`[data-testid="upgrade-button"]`,
		`[data-testid="premium-link"]`,
		`[data-testid="download-app-button"]`,
		`[data-testid="install-app-button"]`,
		`[data-encore-id="buttonTertiary"]`,
		`[data-encore-id="buttonSecondary"]`,
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
		`[class*="CommercialBreak"]`,
		`[class*="ad-slot"]`,
		`[class*="AdSlot"]`,
		`[data-testid="now-playing-bar-ad"]`,
		`[class*="Commercial"]`,
		`[class*="AdBanner"]`,
		`[class*="AdLabel"]`,
		`[class*="NowPlayingAd"]`,
		`[class*="UpgradeButton"]`,
		`[class*="upgrade-button"]`,
		`[class*="upsell"]`,
	}
}