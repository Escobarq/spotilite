// core.js — installs the window.Spicetify namespace skeleton.
// Runs first; every later piece attaches methods onto the object created here.
//
// This is a SUBSET of the spicetify-cli API surface.
// Sub-APIs (uri, cosmos, player, keyboard, components) populate their own keys.
(function initSpicetify() {
  if (typeof window === "undefined") return;
  if (window.Spicetify && window.Spicetify.__spt_ready) return;

  var eventRegistry = (function () {
    var map = Object.create(null);
    return {
      on: function (evt, cb) {
        if (!map[evt]) map[evt] = [];
        map[evt].push(cb);
      },
      off: function (evt, cb) {
        if (!map[evt]) return;
        map[evt] = map[evt].filter(function (x) { return x !== cb; });
      },
      dispatch: function (evt, payload) {
        if (!map[evt]) return;
        for (var i = 0; i < map[evt].length; i++) {
          try { map[evt][i](payload); } catch (e) { console.error("[Spotilite] event", evt, e); }
        }
      }
    };
  })();

  function emitter() {
    var subs = [];
    return {
      subscribe: function (cb) {
        subs.push(cb);
        return function () {
          subs = subs.filter(function (x) { return x !== cb; });
        };
      },
      emit: function (v, ev) {
        for (var i = 0; i < subs.length; i++) {
          try { subs[i](v, ev); } catch (e) { console.error("[Spotilite] sub", e); }
        }
      },
      size: function () { return subs.length; }
    };
  }

  var S = {
    __spt_ready: true,
    __spt_variant: "web",
    version: "0.1.0-web",
    Platform: {
      PlaylistAPI: {},
      LibraryAPI: {},
      PlaybackAPI: {},
      SearchAPI: {},
      HistoryAPI: {},
      TransparencyAPI: {},
      PlaylistPermissionsAPI: {},
      CosmosAPI: {},
      AppSearchAPI: {},
      AuthorizeStatefulAPI: {},
      RemoteConfigResolverAPI: {},
      ShowAPIFilterV2API: {},
      SelectionStateManagerAPI: {},
      LocalStorageAPI: {},
      ConnectivityAPI: {},
      AppStateAPI: {},
      TrackAPI: {},
      RegistrationAPI: {},
      operatingSystem: (navigator.userAgentData && navigator.userAgentData.platform) || navigator.platform || "web",
      PlatformData: { client_version: "web", os: "web" }
    },
    Config: {},
    Player: {
      data: null,
      origin: null,
      _emitter: emitter(),
      addEventListener: function (evt, cb) {
        if (!Spicetify.Player._emitter.subscribe || arguments.length < 2) return undefined;
        return Spicetify.Player._emitter.subscribe(cb);
      },
      dispatchEvent: function (evt) { Spicetify.Player._emitter.emit(undefined, evt && evt.type); },
      _emit: function (type, payload) { Spicetify.Player._emitter.emit(payload, type); }
    },
    LocalStorage: {
      get: function (k) { try { return localStorage.getItem(k); } catch (e) { return null; } },
      set: function (k, v) { try { localStorage.setItem(k, String(v)); } catch (e) {} },
      remove: function (k) { try { localStorage.removeItem(k); } catch (e) {} }
    },
    URI: null,
    CosmosAsync: null,
    Keyboard: null,
    Menu: null,
    ContextMenu: null,
    Topbar: null,
    Playbar: null,
    PopupModal: null,
    ReactComponent: null,
    Panel: null,
    Events: eventRegistry,
    showNotification: function (msg) {
      try {
        if (window.__origNotification) {
          new window.__origNotification(msg).show();
        } else if (typeof Notification !== "undefined" && Notification.permission === "granted") {
          new Notification(msg);
        } else {
          console.info("[Spotilite notify]", msg);
        }
      } catch (e) { console.error("[Spotilite] notify", e); }
    },
    addToQueue: function (uri) { console.warn("[Spotilite] addToQueue not implemented for web player, ignored", uri); },
    removeFromQueue: function (uri) { console.warn("[Spotilite] removeFromQueue not implemented for web player, ignored", uri); },
    colorExtractor: function (_url, _cb) { if (_cb) _cb(["#121212", "#1db954", "#ffffff"]); },
    getAudioData: function () { return null; },
    getFontStyle: function () { return "" },
    Queue: { nextTracks: [], prevTracks: [] },
    SVGIcon: null,
    React: null,
    ReactDOM: null,
    Tippy: null,
    Mousetrap: null,
    Snackbar: null,
    Locale: { getLocale: function () { return navigator.language || "en"; }, getString: function (k) { return k; } }
  };

  window.Spicetify = S;
  window.__spotilite = { spicetify: S, loadedAt: Date.now() };

  // Stub: GraphQL.Definitions is INTENTIONALLY empty on the web variant.
  // Spicetify.cli's preprocess.go populates these by matching minified
  // strings in xpui.js which are unreachable on https://open.spotify.com
  // because the bundle is served by Spotify and we cannot rewrite it.
  // Extensions that rely on specific query names will fail; this is
  // documented in README.
  S.GraphQL = { Definitions: {}, Create: { /* deprecated */ } };

  // Convenience: poll until uri is loaded, then dispatch platformLoaded.
  (function platformWait() {
    var tries = 0;
    var max = 50; // 5s @ 100ms
    var handle = setInterval(function () {
      tries++;
      if (window.Spicetify && window.Spicetify.URI && window.Spicetify.CosmosAsync) {
        clearInterval(handle);
        try { eventRegistry.dispatch("platformLoaded", {}); } catch (e) {}
        return;
      }
      if (tries >= max) {
        clearInterval(handle);
        console.warn("[Spotilite] Spicetify namespace incomplete after wait; some APIs may be undefined");
      }
    }, 100);
  })();
})();
