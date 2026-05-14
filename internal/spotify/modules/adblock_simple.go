package modules

const adblockSimpleCSS = `
[data-testid="ad-type-banner"],
[data-testid="billboard-ad"],
[data-testid="leaderboard-ad"],
[data-testid="sponsorship-ad"],
[data-testid="hpto-ad"],
[data-testid="ad-card"],
[data-testid="ad-slot"],
[data-testid="commercial-break"],
[data-testid="ad-overlay"],
[aria-label="Advertisement"],
[aria-label="Publicidad"],
iframe[src*="doubleclick.net"],
iframe[src*="moatads"] {
display: none !important;
}
`

const adblockSimpleJS = `
(function() {
if (window.__spotiliteAdblockInstalled) return;
window.__spotiliteAdblockInstalled = true;

var LOCAL_API = 'localhost:8765';

var AD_URL_PATTERNS = [
'doubleclick.net',
'googlesyndication.com',
'googleadservices.com',
'moatads.com',
'adservice.google.com',
'audio-ads.spotify.com',
'ads-audio.spotify.com',
'ad-logger.spotify.com',
'ad-handler.spotify.com',
'spotify.com/ads/',
'spotify.com/audioad',
'audio-sp.spotify.com/ad',
'audio-fa.scdn.co/ad',
'heads-fa.scdn.co/ad',
'partnerakamai.spotify.com/ad',
'megaphone.fm/ad',
'/ad/',
'/ads/',
'adclick',
'adserver'
];

function isAdUrl(url) {
if (!url) return false;
var u = typeof url === 'string' ? url : (url.url ? url.url : '');
if (!u || u.indexOf(LOCAL_API) !== -1) return false;
if (u.indexOf('spotify.com') !== -1 && u.indexOf('/ad') === -1 && u.indexOf('/ads') === -1) return false;
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

console.log('[Spotilite AdBlock] Simple network-only mode');
console.log('[Spotilite AdBlock] For better ad blocking, use SpotX: https://github.com/SpotX-Official/SpotX');
})();
`

type AdBlockSimpleModule struct {
	BaseModule
}

func NewAdBlockSimpleModule(enabled bool) *AdBlockSimpleModule {
	return &AdBlockSimpleModule{
		BaseModule: BaseModule{name: "adblock_simple", enabled: enabled},
	}
}

func (m *AdBlockSimpleModule) CSS() string { return adblockSimpleCSS }
func (m *AdBlockSimpleModule) JS() string  { return adblockSimpleJS }
func (m *AdBlockSimpleModule) Selectors() []string {
	return []string{
		`[data-testid="ad-type-banner"]`,
		`[data-testid="billboard-ad"]`,
		`[data-testid="commercial-break"]`,
		`[aria-label="Advertisement"]`,
		`[aria-label="Publicidad"]`,
	}
}
