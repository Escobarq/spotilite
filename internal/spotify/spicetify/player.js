// player.js — Spicetify.Player reimplementation.
// The desktop client leaks the internal Redux/Spotify.Player state through the
// regex-patched xpui.js. The web player does not expose that object, so we
// rebuild a similar facade from:
//
//   - navigator.mediaSession.metadata (track title/artist/album/source URL)
//   - DOM polling on .now-playing-bar for play/pause/shuffle/repeat state
//   - simulated clicks on Spotify's own control buttons for commands
//
// The contract exposed matches the most-used subset of spicetify.cli's Player:
// addEventListener, getProgress, seek, play, pause, next, back, etc.
(function initSpicetifyPlayer() {
  if (!window.Spicetify) return;

  function $(sel, root) { return (root || document).querySelector(sel); }
  function $$(sel, root) {
    return Array.prototype.slice.call((root || document).querySelectorAll(sel));
  }

  function buttonByLabel(labels) {
    var lower = labels.map(function (l) { return l.toLowerCase(); });
    return $$("button").filter(function (b) {
      var al = (b.getAttribute("aria-label") || "").toLowerCase();
      return lower.indexOf(al) !== -1;
    });
  }

  function electronClick(el) {
    if (!el) return;
    try { el.click(); } catch (e) {}
  }

  function dispatchClick(el) {
    if (!el) return;
    var ev = new MouseEvent("click", { bubbles: true, cancelable: true });
    el.dispatchEvent(ev);
  }

  // State mirrors spicetify.cli's Player._state shape loosely.
  var state = {
    isPaused: null,
    track: null,
    positionMs: 0,
    durationMs: 0,
    repeat: 0,
    shuffle: false,
    context: null
  };

  function readMeta() {
    var meta = navigator.mediaSession && navigator.mediaSession.metadata;
    var t = null;
    if (meta) {
      t = {
        uri: null,
        url: null,
        title: meta.title || "",
        artists: meta.artist ? meta.artist.split(",") : [],
        album: meta.album || ""
      };
      // URI not present in MediaSession; we tag it from a hidden data-* if any.
      var np = document.querySelector(".main-nowPlayingBar-nowPlayingBar");
      if (np) {
        var uriEl = np.querySelector("[data-testid='track-info'] a[href*='/track/']");
        if (uriEl) {
          var href = uriEl.getAttribute("href");
          var m = /\/track\/([a-zA-Z0-9]+)/.exec(href);
          if (m) t.uri = "spotify:track:" + m[1];
        }
      }
    }
    return t;
  }

  function readPauseState() {
    var pauseBtn = buttonByLabel(["Pause"])[0];
    if (pauseBtn) {
      return { isPaused: false, label: "Pause" };
    }
    var playBtn = buttonByLabel(["Play"])[0];
    if (playBtn) {
      return { isPaused: true, label: "Play" };
    }
    return { isPaused: null, label: null };
  }

  function poll() {
    try {
      var pause = readPauseState();
      var prevPaused = state.isPaused;
      state.isPaused = pause.isPaused;
      if (prevPaused !== pause.isPaused && pause.isPaused !== null) {
        Spicetify.Player._emit(pause.isPaused ? "playpause" : "playpause");
      }

      var meta = readMeta();
      if (meta && state.track !== meta && meta.uri) {
        state.track = meta;
        Spicetify.Player._emit("songchange", meta);
      }

      // duration via media element
      var audio = document.querySelector("audio");
      if (audio && !isNaN(audio.duration * 1000)) {
        state.durationMs = Math.round(audio.duration * 1000);
        state.positionMs = Math.round((audio.currentTime || 0) * 1000);
      }
    } catch (e) { /* ignore poll errors */ }
  }

  setInterval(poll, 1000);
  // Run an initial poll on next tick so the wrapper is queryable immediately.
  setTimeout(poll, 50);

  var Player = {
    data: state,
    origin: state,
    addEventListener: function (type, cb) {
      // Simulates a DOM EventTarget addEventListener("songchange", cb).
      // subscribe returns an unsubscribe function.
      var prev = Spicetify.Player._emitter.size();
      var unsub = Spicetify.Player._emitter.subscribe(function (v, evt) {
        if (evt === type) {
          try { cb(Object.assign({ type: type }, v || {})); } catch (e) {}
        }
      });
      return unsub;
    },
    removeEventListener: function (type, cb) {
      // Spicetify.cli stores the same emitter; for simplicity we don't track
      // separate callbacks. accept and ignore. Most extensions never unregister.
    },
    dispatchEvent: function (e) { Spicetify.Player._emitter.emit(e, e && e.type); },
    play: function () {
      var btn = buttonByLabel(["Play"])[0];
      if (btn) dispatchClick(btn);
      Spicetify.Player._emit("play");
    },
    pause: function () {
      var btn = buttonByLabel(["Pause"])[0];
      if (btn) dispatchClick(btn);
      Spicetify.Player._emit("pause");
    },
    togglePlay: function () {
      if (state.isPaused) Spicetify.Player.play(); else Spicetify.Player.pause();
    },
    next: function () {
      var btn = buttonByLabel(["Next"])[0];
      if (btn) dispatchClick(btn);
      Spicetify.Player._emit("next");
    },
    back: function () {
      var btn = buttonByLabel(["Previous"])[0];
      if (btn) dispatchClick(btn);
      Spicetify.Player._emit("back");
    },
    seek: function (ms) {
      var audio = document.querySelector("audio");
      if (audio) audio.currentTime = (ms || 0) / 1000;
      Spicetify.Player._emit("seek", ms);
    },
    seekMs: function (ms) { return Spicetify.Player.seek(ms); },
    getProgress: function () { return { position: state.positionMs | 0, duration: state.durationMs | 0 }; },
    getCurrentState: function () {
      var t = state.track || {};
      return {
        track: t,
        position: state.positionMs,
        duration: state.durationMs,
        isPaused: state.isPaused,
        isBuffering: false,
        isRepeating: !!state.repeat,
        repeatMode: state.repeat,
        isShuffling: !!state.shuffle,
        context: state.context,
        previousTracks: Spicetify.Queue.prevTracks.slice(-50),
        nextTracks: Spicetify.Queue.nextTracks.slice(0, 50)
      };
    },
    toggleShuffle: function () {
      state.shuffle = !state.shuffle;
      var btn = buttonByLabel(["Shuffle"])[0];
      if (btn && state.shuffle === false) dispatchClick(btn);
      // The shuffle button on web toggles on click, so we do not force the
      // second click here once it's on; just rely on state.
      Spicetify.Player._emit("shuffle", state.shuffle);
    },
    toggleRepeat: function () {
      state.repeat = (state.repeat + 1) % 3;
      var btn = buttonByLabel(["Repeat"])[0];
      if (btn) dispatchClick(btn);
      Spicetify.Player._emit("repeat", state.repeat);
    },
    setVolume: function (vol) {
      var audio = document.querySelector("audio");
      if (audio) audio.volume = Math.max(0, Math.min(1, vol));
    },
    mute: function () { Spicetify.Player.setVolume(0); },
    unmute: function () { Spicetify.Player.setVolume(1); },
    toggleMute: function () {
      var audio = document.querySelector("audio");
      if (audio) audio.muted = !audio.muted;
    },
    toggleHeart: function () {
      var btn = buttonByLabel(["Save to your Liked Songs", "Remove from your Liked Songs"])[0];
      if (btn) dispatchClick(btn);
    },
    formatTime: function (seconds) {
      var m = Math.floor((seconds || 0) / 60);
      var s = Math.floor((seconds || 0) % 60);
      return m + ":" + (s < 10 ? "0" + s : s);
    }
  };

  window.Spicetify.Player = Player;
})();
