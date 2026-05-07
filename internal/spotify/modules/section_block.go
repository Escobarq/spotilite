package modules

const sectionBlockJS = `
(function() {
	if (window.__spotiliteSectionBlockInstalled) return;
	window.__spotiliteSectionBlockInstalled = true;

	var BLOCKED_SECTIONS = {
		'0JQ5DAnM3wGh0gz1MXnul1': 'Party',
		'0JQ5DAnM3wGh0gz1MXnukV': 'Chill',
		'0JQ5IMCbQBLupUQrQFeCzx': 'Best of the Year',
		'0JQ5DAnM3wGh0gz1MXnu3C': 'Best of Artists / Tracks',
		'0JQ5DAnM3wGh0gz1MXnu4w': 'Best of songwriters',
		'0JQ5IMCbQBLhSb02SGYpDM': 'Biggest Indie Playlists',
		'0JQ5DAnM3wGh0gz1MXnu5g': 'Charts',
		'0JQ5DAnM3wGh0gz1MXnu3p': 'Dinner',
		'0JQ5DAob0KOew1FBAMSmBz': 'Featured Charts',
		'0JQ5DAob0JCuWaGLU6ntFY': 'Focus',
		'0JQ5DAnM3wGh0gz1MXnulP': 'Focus',
		'0JQ5DAnM3wGh0gz1MXnu3s': 'Fresh new music',
		'0JQ5DAob0LaV9FOMJ9utY5': 'Gaming music',
		'0JQ5DAnM3wGh0gz1MXnu3q': 'Happy',
		'0JQ5IMCbQBLiqrNCH9VvmA': 'ICE PHONK',
		'0JQ5DAnM3wGh0gz1MXnucG': 'Mood',
		'0JQ5DAob0JCuWaGLU6ntFT': 'Mood',
		'0JQ5IMCbQBLicmNERjnGn5': 'Most Listened 2023',
		'0JQ5DAob0Jr9ClCbkV4pZD': 'Music to game to',
		'0JQ5DAnM3wGh0gz1MXnu3B': 'Popular Albums / Artists',
		'0JQ5DAnM3wGh0gz1MXnu3D': 'Popular new releases',
		'0JQ5DAnM3wGh0gz1MXnu4h': 'Popular radio',
		'0JQ5DAnM3wGh0gz1MXnu3u': 'Sad',
		'0JQ5DAnM3wGh0gz1MXnul2': 'Sad',
		'0JQ5DAnM3wGh0gz1MXnu3w': 'Throwback',
		'0JQ5DAnM3wGh0gz1MXnul4': 'Throwback',
		'0JQ5DAuChZYPe9iDhh2mJz': 'Throwback Thursday',
		'0JQ5DAnM3wGh0gz1MXnu3M': 'Todays biggest hits',
		'0JQ5DAnM3wGh0gz1MXnu3E': 'Trending now',
		'0JQ5DAnM3wGh0gz1MXnu3x': 'Workout',
		'0JQ5DAnM3wGh0gz1MXnul6': 'Workout',
		'0JQ5IMCbQBLlC31GvtaB6w': 'Now defrosting',
		'0JQ5IMCbQBLqTJyy28YCa9': 'Unknown',
		'0JQ5DAnM3wGh0gz1MXnu7R': 'Unknown'
	};

	var BLOCKED_CONTENT_TYPES = new Set(['Podcast', 'Audiobook', 'Episode']);

	var API_PATHFINDER = 'api-partner.spotify.com/pathfinder';
	var API_RECOMMENDATIONS = 'api.spotify.com/v1/views/personalized-recommendations';

	var _prevFetch = window.fetch;

	function createAdapter(isPR) {
		if (isPR) {
			return {
				getId: function(item) {
					var href = item && item.href;
					if (!href) return null;
					var parts = href.split('/');
					var id = parts[parts.length - 1];
					if (id && id.startsWith('section')) id = id.substring(7);
					return id;
				},
				getTitle: function(item) { return (item && item.content && item.content.name) || 'Unknown'; },
				getContentItems: function(item) { return item && item.content && item.content.items; },
				getContentData: function(ci) { return ci && ci.content; },
				getContentType: function(ci) { return ci && ci.type; },
				getContentTypeName: function(ci) { return ci && ci.content_type; },
				getSectionId: function(item) { return item && item.id; }
			};
		}
		return {
			getId: function(item) {
				var uri = item && item.uri;
				if (!uri) return null;
				var parts = uri.split(':');
				return parts[parts.length - 1];
			},
			getTitle: function(item) { return (item && item.data && item.data.title && item.data.title.text) || 'Unknown'; },
			getContentItems: function(item) { return item && item.sectionItems && item.sectionItems.items; },
			getContentData: function(ci) { return ci && ci.content && ci.content.data; },
			getContentType: function() { return null; },
			getContentTypeName: function() { return null; },
			getSectionId: function() { return null; }
		};
	}

	function removeSections(items, adapter) {
		if (!items || !items.length) return;
		var removed = [];
		for (var i = items.length - 1; i >= 0; i--) {
			var sectionId = adapter.getId(items[i]);
			if (sectionId && BLOCKED_SECTIONS[sectionId]) {
				removed.push({ id: sectionId, name: BLOCKED_SECTIONS[sectionId] });
				items.splice(i, 1);
			}
		}
		if (removed.length > 0) console.log('[Spotilite] Removed ' + removed.length + ' blocked section(s)');
	}

	function removePodcasts(items, adapter, isPR) {
		if (!items || !items.length) return;
		var removed = [];
		for (var i = items.length - 1; i >= 0; i--) {
			var contentItems = adapter.getContentItems(items[i]);
			if (!contentItems || !contentItems.length) continue;

			if (isPR) {
				var sectionId = adapter.getSectionId(items[i]);
				if (sectionId === 'shortcuts') {
					for (var j = contentItems.length - 1; j >= 0; j--) {
						var ct = adapter.getContentTypeName(contentItems[j]);
						if (ct === 'PODCAST_EPISODE' || ct === 'AUDIOBOOK') {
							contentItems.splice(j, 1);
						}
					}
					continue;
				}
				var firstType = adapter.getContentType(contentItems[0]);
				if (firstType === 'show') {
					removed.push(adapter.getTitle(items[i]));
					items.splice(i, 1);
					continue;
				}
			}

			for (var k = contentItems.length - 1; k >= 0; k--) {
				var data = adapter.getContentData(contentItems[k]);
				if (data && BLOCKED_CONTENT_TYPES.has(data.__typename)) {
					contentItems.splice(k, 1);
				}
			}
		}
		if (removed.length > 0) console.log('[Spotilite] Removed ' + removed.length + ' podcast section(s)');
	}

	function removeCanvasSections(sections) {
		if (!sections || !sections.length) return;
		for (var i = sections.length - 1; i >= 0; i--) {
			if (sections[i] && sections[i].data && sections[i].data.__typename === 'HomeFeedBaselineSectionData') {
				sections.splice(i, 1);
			}
		}
	}

	function processHomeData(data) {
		var body = data && data.data && data.data.home;
		var sections = body && body.sectionContainer && body.sectionContainer.sections && body.sectionContainer.sections.items;
		var items = (data && data.content && data.content.items) || (data && data.data && data.data.content && data.data.content.items);
		var isPR = !!items && !body;
		var targetArray = isPR ? items : sections;

		if (!targetArray || !targetArray.length) return;

		var adapter = createAdapter(isPR);
		removeSections(targetArray, adapter);
		removePodcasts(targetArray, adapter, isPR);
		if (!isPR && sections) removeCanvasSections(sections);
	}

	window.fetch = function() {
		var args = arguments;
		var url = args[0];
		var urlString = typeof url === 'string' ? url : (url && url.url ? url.url : '');

		var isPF = urlString.indexOf(API_PATHFINDER) !== -1;
		var isPR = urlString.indexOf(API_RECOMMENDATIONS) !== -1;

if (!isPF && !isPR) return _prevFetch.apply(this, args);

return _prevFetch.apply(this, args).then(function(response) {
			var cloned = response.clone();
			return response.json().then(function(data) {
				var shouldModify = (isPF && data && data.data && data.data.home) || (isPR && data && data.content);
				if (!shouldModify) return cloned;

				processHomeData(data);

				return new Response(JSON.stringify(data), {
					status: response.status,
					statusText: response.statusText,
					headers: response.headers
				});
			}).catch(function() { return cloned; });
		});
	};
})();
`

type SectionBlockModule struct {
	BaseModule
}

func NewSectionBlockModule(enabled bool) *SectionBlockModule {
	return &SectionBlockModule{
		BaseModule: BaseModule{name: "sectionblock", enabled: enabled},
	}
}

func (m *SectionBlockModule) JS() string { return sectionBlockJS }
