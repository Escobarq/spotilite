package modules

const historyJS = `
(function() {
	if (window.__spotiliteHistoryInstalled) return;
	window.__spotiliteHistoryInstalled = true;

	var MAX_TRACKS = 1000;
	var MAX_DELAY = 1000;

	var STORAGE_KEY = 'spotilite.sentTracks';

	function loadTracks() {
		try {
			var saved = localStorage.getItem(STORAGE_KEY);
			return saved ? new Set(JSON.parse(saved)) : new Set();
		} catch(e) { return new Set(); }
	}

	function saveTracks(tracks) {
		try {
			var arr = Array.from(tracks);
			if (arr.length > MAX_TRACKS) arr = arr.slice(-MAX_TRACKS);
			localStorage.setItem(STORAGE_KEY, JSON.stringify(arr));
		} catch(e) {}
	}

	var unique = loadTracks();
	var timeout;

	function debounce(fn, wait) {
		return function() {
			var args = arguments;
			var later = function() { clearTimeout(timeout); fn.apply(null, args); };
			clearTimeout(timeout);
			timeout = setTimeout(later, wait);
		};
	}

	var _origPostMessage = window.postMessage;
	window.addEventListener('message', function(e) {
		if (!e || !e.data) return;
		var d = e.data;
		if (d.type === 'playback-state' && d.item && d.item.uri) {
			debouncedHandler(d);
		}
	});

	var debouncedHandler = debounce(function(data) {
		var uri = data.item.uri;
		if (uri && uri.indexOf('spotify:track:') !== -1 && !unique.has(uri)) {
			unique.add(uri);
			saveTracks(unique);
			console.log('[Spotilite] Track logged: ' + uri);
		}
	}, MAX_DELAY);

var _prevFetch = window.fetch;
var lastTrackUri = '';
window.fetch = function() {
var args = arguments;
return _prevFetch.apply(this, args).then(function(response) {
			try {
				var url = args[0];
				if (typeof url === 'string' && url.indexOf('spclient.wg.spotify.com/player') !== -1) {
					var cloned = response.clone();
					cloned.json().then(function(data) {
						if (data && data.item && data.item.uri && data.item.uri.indexOf('spotify:track:') !== -1) {
							if (lastTrackUri !== data.item.uri) {
								lastTrackUri = data.item.uri;
								if (!unique.has(data.item.uri)) {
									unique.add(data.item.uri);
									saveTracks(unique);
								}
							}
						}
					}).catch(function(){});
				}
			} catch(e) {}
			return response;
		});
	};
})();
`

type HistoryModule struct {
	BaseModule
}

func NewHistoryModule(enabled bool) *HistoryModule {
	return &HistoryModule{
		BaseModule: BaseModule{name: "history", enabled: enabled},
	}
}

func (m *HistoryModule) JS() string { return historyJS }
