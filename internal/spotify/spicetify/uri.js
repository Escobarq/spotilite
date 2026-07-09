// uri.js — minimal spicetify URI parser/serializer compatible with the subset
// used by most community extensions. Self-contained (no WeakRef/Symbol) so it
// works inside the strict CSP webview context.
//
// This is a simplified port: parses spotify:track/.../album/.../playlist/...
// artist/.../episode/.../show/.../user/.../search/... and exposes Parse/Object/
// toURI helpers similar to spicetify.cli's URI.js.
(function initSpicetifyURI() {
  if (!window.Spicetify) return;

  function URI(s) {
    if (!(this instanceof URI)) return new URI(s);
    this._raw = s;
    var parsed = URI.parse(s);
    this.type = parsed.type;
    this._base62Id = parsed.id;
    this._category = parsed.category;
    this._arguments = parsed.arguments;
  }

  URI.prototype.id = function () {
    if (this._base62Id) return this._base62Id;
    if (this._arguments && this._arguments.length && this._arguments[0]) return this._arguments[0];
    return null;
  };

  URI.prototype.toString = function () {
    if (this._raw) return this._raw;
    var parts = ["spotify", this.type];
    if (this._category) parts.push(this._category);
    if (this._base62Id) parts.push(this._base62Id);
    if (this._arguments && this._arguments.length) {
      for (var i = 0; i < this._arguments.length; i++) parts.push(this._arguments[i]);
    }
    return parts.join(":");
  };

  URI.prototype.toURI = function () { return new URI(this.toString()); };

  URI.prototype.equals = function (other) {
    if (!other) return false;
    return this.toString() === (other.toString ? other.toString() : String(other));
  };

  URI.prototype.isEmpty = function () { return !this._raw && !this.type; };

  URI.Type = Object.freeze({
    AD: "ad", ALBUM: "album", ARTIST: "artist", EPISODE: "episode",
    SHOW: "show", TRACK: "track", PLAYLIST: "playlist", USER: "user",
    SEARCH: "search", RADIO: "radio", CONCERT: "concert",
    APPLICATION: "application", UNKNOWN: "unknown"
  });

  var PREFIX_RE = /^spotify:(track|album|artist|playlist|episode|show|user|search|ad|concert|application|radio|local|prereleasepolicy)(?::|$)/;

  function classify(type) {
    var t = (type || "").toLowerCase();
    if (t === "track") return URI.Type.TRACK;
    if (t === "album") return URI.Type.ALBUM;
    if (t === "artist") return URI.Type.ARTIST;
    if (t === "playlist") return URI.Type.PLAYLIST;
    if (t === "episode") return URI.Type.EPISODE;
    if (t === "show") return URI.Type.SHOW;
    if (t === "user") return URI.Type.USER;
    if (t === "search") return URI.Type.SEARCH;
    if (t === "ad") return URI.Type.AD;
    if (t === "concert") return URI.Type.CONCERT;
    if (t === "application") return URI.Type.APPLICATION;
    if (t === "radio") return URI.Type.RADIO;
    if (t === "prereleasepolicy") return URI.Type.PRERELEASEPOLICY;
    return URI.Type.UNKNOWN;
  }

  URI.parse = function (s) {
    if (!s || typeof s !== "string") {
      return { type: URI.Type.UNKNOWN, id: null, arguments: [] };
    }
    var parts = s.split(":");
    if (parts[0] !== "spotify") {
      return { type: URI.Type.UNKNOWN, id: null, arguments: parts.slice(1) };
    }
    var m = PREFIX_RE.exec(s);
    if (!m) {
      return { type: URI.Type.UNKNOWN, id: null, arguments: parts.slice(2) };
    }
    var type = classify(m[1]);
    var rest = parts.slice(2);
    var id = null;
    var category = null;
    if (type === URI.Type.PLAYLIST || type === URI.Type.ALBUM || type === URI.Type.ARTIST || type === URI.Type.SHOW || type === URI.Type.EPISODE) {
      // spotify:playlist:<base62-id> or spotify:playlist:user:<user-id>:...
      // spotify:album:n5NX...(63-char)  or spotify:show:abc:... or episode:...
      if (/^[a-zA-Z0-9]{16,40}$/.test(rest[0])) {
        id = rest[0];
      } else {
        // could be playlist:user:<user>:base62:...
        if (rest.length > 1 && /^[a-zA-Z0-9]+$/.test(rest[1])) {
          category = rest[0];
          id = rest[1];
          rest = rest.slice(2);
        } else {
          category = rest[0];
          rest = rest.slice(1);
          id = null;
        }
      }
    } else if (type === URI.Type.USER) {
      id = rest[0] || null;
      rest = rest.slice(1);
    } else {
      id = rest[0] || null;
      rest = rest.slice(1);
    }
    return { type: type, id: id, category: category, arguments: rest };
  };

  URI.isValid = function (s) {
    if (!s || typeof s !== "string") return false;
    return PREFIX_RE.test(s);
  };

  URI.from = URI.parse;

  // Attach to namespace
  window.Spicetify.URI = URI;
})();
