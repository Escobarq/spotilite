// cosmos.js — Spicetify.CosmosAsync reimplementation.
// The desktop client routes Cosmos traffic through internal wg:// and sp://
// schemes registered by Spotify's CEF bridge. On https://open.spotify.com
// those schemes don't exist; instead we route to the public web API:
//
//   https://api.spotify.com/v1/...
//   https://spclient.wg.spotify.com/...
//   https://open.spotify.com/...
//
// The Bearer token is read from localStorage keys that the web player's
// session bootstraps, with fallbacks. Requests are first attempted directly;
// if a CORS preflight fails, the request is retried through the local Go
// proxy at http://localhost:8765/proxy?url=<encoded>.
//
// Extensions call Spicetify.CosmosAsync.get/post/request(method, url, body, headers).
// The interface is intentionally similar enough to spicetify.cli's that
// common extensions work unchanged for endpoints exposed by the web API.
(function initSpicetifyCosmos() {
  if (!window.Spicetify) return;

  // Look for auth token candidates that the Spotify web player keeps around.
  // The actual key is set by Spotify's bootstrap; we probe a few common keys and
  // fall back to reading sessionStorage and a "spotify-token" cookie.
  function readToken() {
    var tries = [
      function () { return localStorage.getItem("wp-token-async"); },
      function () { return localStorage.getItem("spotify-token"); },
      function () { return localStorage.getItem("accessToken"); },
      function () {
        var c = document.cookie.split("; ");
        for (var i = 0; i < c.length; i++) {
          var p = c[i].split("=");
          if (p[0] === "spotify-token") return decodeURIComponent(p.slice(1).join("="));
          if (p[0] === "accessToken") return decodeURIComponent(p.slice(1).join("="));
        }
        return null;
      }
    ];
    for (var i = 0; i < tries.length; i++) {
      try {
        var v = tries[i]();
        if (v) return v.replace(/^"|"$/g, "");
      } catch (e) { /* ignore */ }
    }
    return null;
  }

  function getBaseOrigin() {
    var m = /^(https?:\/\/[^\/]+)/.exec(location.href);
    return m ? m[1] : "https://open.spotify.com";
  }

  function buildHeaders(extra) {
    var h = {
      "App-Platform": "WebPlayer",
      "Spotify-App-Version": "Web",
      "Accept": "application/json",
      "Origin": getBaseOrigin(),
      "Referer": getBaseOrigin() + "/"
    };
    var token = readToken();
    if (token) h["Authorization"] = "Bearer " + token;
    if (extra) {
      for (var k in extra) if (Object.prototype.hasOwnProperty.call(extra, k)) h[k] = extra[k];
    }
    return h;
  }

  function normalizeUrl(url) {
    if (!url) return url;
    if (url.indexOf("://") === 0) {
      url = "https" + url; // schemeless
    }
    // The desktop sp:// and wg:// schemes map to web endpoints if any.
    url = url.replace(/^sp:\/\//, "https://open.spotify.com/");
    url = url.replace(/^wg:\/\/[^\/]*\//, "https://spclient.wg.spotify.com/");
    url = url.replace(/^hm:\/\/[^\/]*\//, "https://api.spotify.com/");
    return url;
  }

  function directFetch(url, init) {
    url = normalizeUrl(url);
    return fetch(url, init);
  }

  function proxiedFetch(url, init) {
    url = normalizeUrl(url);
    var proxy = (window.__spotiliteApiBase || "http://localhost:8765") + "/proxy";
    var proxyUrl = proxy + "?url=" + encodeURIComponent(url);
    return fetch(proxyUrl, init);
  }

  function attempt(url, init, allowProxy) {
    return new Promise(function (resolve, reject) {
      directFetch(url, init).then(function (resp) {
        // https://open.spotify.com's API often responds with 401 when the local
        // token is stale but the proxy can refresh it. Retry through the proxy
        // if we are allowed and the call looks like a Spotify endpoint.
        if (resp.status === 401 && allowProxy && /^https:\/\/(api\.spotify\.com|spclient\.wg\.spotify\.com|open\.spotify\.com)/.test(url)) {
          proxiedFetch(url, init).then(resolve, reject);
          return;
        }
        resolve(resp);
      }, function (err) {
        if (allowProxy) proxiedFetch(url, init).then(resolve, reject);
        else reject(err);
      });
    });
  }

  function readJson(resp) {
    return new Promise(function (resolve, reject) {
      if (resp.status === 204) return resolve(null);
      var ct = (resp.headers.get("content-type") || "").toLowerCase();
      if (ct.indexOf("json") === -1) return resolve(resp);
      resp.clone().text().then(function (t) {
        if (!t) return resolve(null);
        try { resolve(JSON.parse(t)); } catch (e) { resolve(t); }
      }, reject);
    });
  }

  function Cosmos(method, url, body, headers, options) {
    options = options || {};
    var init = { method: method, headers: buildHeaders(headers), credentials: "same-origin" };
    if (body !== undefined && body !== null) {
      if (typeof body === "string" || body instanceof Blob || body instanceof ArrayBuffer) {
        init.body = body;
      } else {
        init.body = JSON.stringify(body);
        if (!init.headers["Content-Type"]) init.headers["Content-Type"] = "application/json";
      }
    }
    return attempt(url, init, !options.noproxy).then(readJson);
  }

  // Public surface
  window.Spicetify.CosmosAsync = {
    get:    function (url, headers, options) { return Cosmos("GET", url, null, headers, options); },
    post:   function (url, body, headers, options) { return Cosmos("POST", url, body, headers, options); },
    put:    function (url, body, headers, options) { return Cosmos("PUT", url, body, headers, options); },
    del:    function (url, headers, options) { return Cosmos("DELETE", url, null, headers, options); },
    patch:  function (url, body, headers, options) { return Cosmos("PATCH", url, body, headers, options); },
    head:   function (url, headers, options) { return Cosmos("HEAD", url, null, headers, options); },
    request: function (method, url, body, headers, options) { return Cosmos(method.toUpperCase(), url, body, headers, options); },
    // internal helper exposed for debug
    _readToken: readToken,
    _normalize: normalizeUrl
  };
})();
