// keyboard.js — minimal Spicetify.Keyboard reimplementation.
//
// Goals:
//   - Spicetify.Keyboard.registerKeys(map) -> unsubscribe
//   - Spicetify.Keyboard.registerKey(key, cb) -> unsubscribe
//   - Spicetify.Keyboard.unregisterKey(key)
//   - Spicetify.Keyboard.pressKey(key)   (synthesize)
//
// Keyboard combos follow the spicetify CLI friendly subscript form:
//
//   "ctrl+s", "shift+k", "alt+z", "meta+a", or single characters / digits /
//   function keys "f1"..".
//
// We install a `keydown` listener on window; ignore events coming from
// input/textarea/contenteditable unless the combo bypasses by being a
// modifier-only combo.
(function initSpicetifyKeyboard() {
  if (!window.Spicetify) return;

  function normalize(key) {
    return String(key || "").toLowerCase().replace(/\s+/g, "");
  }

  function isTypingTarget(target) {
    if (!target) return false;
    var t = target.tagName;
    if (t === "INPUT" || t === "TEXTAREA" || t === "SELECT") return true;
    if (target.isContentEditable) return true;
    return false;
  }

  function matchCombo(combo, ev) {
    var parts = combo.split("+");
    var needKey = null;
    var need = { ctrl: false, shift: false, alt: false, meta: false };
    for (var i = 0; i < parts.length; i++) {
      var p = parts[i];
      if (p === "ctrl" || p === "control") need.ctrl = true;
      else if (p === "shift") need.shift = true;
      else if (p === "alt" || p === "option") need.alt = true;
      else if (p === "meta" || p === "cmd" || p === "command") need.meta = true;
      else needKey = p;
    }
    if (need.ctrl !== ev.ctrlKey) return false;
    if (need.shift !== ev.shiftKey) return false;
    if (need.alt !== ev.altKey) return false;
    if (need.meta !== ev.metaKey) return false;
    if (!needKey) return true; // modifier-only combo
    var k = (ev.key || "").toLowerCase();
    if (needKey === k) return true;
    if (needKey === "space" && k === " ") return true;
    if (needKey === "esc" && k === "escape") return true;
    return false;
  }

  var bindings = []; // {key: normalized, handler: fn}

  window.addEventListener("keydown", function (ev) {
    var target = ev.target;
    var inTextField = isTypingTarget(target);
    for (var i = 0; i < bindings.length; i++) {
      var b = bindings[i];
      if (matchCombo(b.key, ev)) {
        // Allow modifiers-only combos while in text fields.
        var isModifierOnly = !/[a-z0-9]/i.test(ev.key) || ev.ctrlKey || ev.metaKey || ev.altKey;
        if (inTextField && !isModifierOnly) continue;
        try {
          var returns = b.handler(ev);
          if (returns === false || (returns && returns.preventDefault !== undefined && returns.preventDefault())) ev.preventDefault();
          else ev.preventDefault();
        } catch (e) { console.error("[Spotilite] Keyboard handler", e); }
        return;
      }
    }
  }, true);

  function registerKey(key, cb) {
    var k = normalize(key);
    var entry = { key: k, handler: cb };
    bindings.push(entry);
    return function () {
      bindings = bindings.filter(function (x) { return x !== entry; });
    };
  }

  function registerKeys(map) {
    // map: { "ctrl+s": cb, "shift+k": cb }
    var keys = Object.keys(map || {});
    var unsubs = [];
    for (var i = 0; i < keys.length; i++) {
      unsubs.push(registerKey(keys[i], map[keys[i]]));
    }
    return function () { for (var i = 0; i < unsubs.length; i++) unsubs[i](); };
  }

  function unregisterKey(key) {
    var k = normalize(key);
    bindings = bindings.filter(function (b) { return b.key !== k; });
  }

  function pressKey(key) {
    var map = { ctrl: false, shift: false, alt: false, meta: false };
    var parts = normalize(key).split("+");
    var code = "";
    for (var i = 0; i < parts.length; i++) {
      var p = parts[i];
      if (p === "ctrl" || p === "control") map.ctrl = true;
      else if (p === "shift") map.shift = true;
      else if (p === "alt" || p === "option") map.alt = true;
      else if (p === "meta" || p === "cmd") map.meta = true;
      else code = p;
    }
    var ev = new KeyboardEvent("keydown", Object.assign({ key: code, bubbles: true, cancelable: true }, map));
    window.dispatchEvent(ev);
  }

  // Spicetify-style API: also expose Mousetrap-like aliases.
  window.Spicetify.Keyboard = {
    registerKey: registerKey,
    registerKeys: registerKeys,
    unregisterKey: unregisterKey,
    pressKey: pressKey,
    reset: function () { bindings = []; }
  };
})();
