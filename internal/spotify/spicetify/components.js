// components.js — minimal Menu / ContextMenu / Topbar / Playbar / PopupModal / Panel
// reimplementations for the web player.
//
// The desktop spicetify-cli exposes Spicetify.ReactComponent.* factories backed
// by Spotify's own React+webpack modules. On the web player we cannot reach
// those internals reliably, so this file provides a DOM-based subset that lets
// common community extensions register buttons and menu items.
//
// Limitations (documented in README):
//   - React-style JSX props are not supported; components mount real DOM nodes.
//   - Spicetify.ReactComponent.ButtonPrimary etc. are factories returning an
//     object with .mount(parent) / .unmount() rather than JSX elements.
//   - Some visually advanced components (Slider, TooltipWrapper) are stubs.
(function initSpicetifyComponents() {
  if (!window.Spicetify) return;

  function el(tag, attrs, children) {
    var node = document.createElement(tag);
    if (attrs) {
      Object.keys(attrs).forEach(function (k) {
        if (k === "className") node.className = attrs[k];
        else if (k === "style" && typeof attrs[k] === "object") Object.assign(node.style, attrs[k]);
        else if (k === "innerHTML") node.innerHTML = attrs[k];
        else if (k === "textContent") node.textContent = attrs[k];
        else if (k.indexOf("on") === 0 && typeof attrs[k] === "function") {
          node.addEventListener(k.slice(2).toLowerCase(), attrs[k]);
        } else node.setAttribute(k, attrs[k]);
      });
    }
    if (children) {
      if (!Array.isArray(children)) children = [children];
      children.forEach(function (c) {
        if (c == null) return;
        node.appendChild(typeof c === "string" ? document.createTextNode(c) : c);
      });
    }
    return node;
  }

  function findSelector(selector) {
    return document.querySelector(selector);
  }

  // ---- Menu / ContextMenu -------------------------------------------------

  function MenuItem(props) {
    if (!(this instanceof MenuItem)) return new MenuItem(props);
    this.props = props || {};
    this._node = null;
  }
  MenuItem.prototype.register = function () {
    var self = this;
    Spicetify.Menu._items.push(self);
    Spicetify.Menu._render();
    return function () { self.unregister(); };
  };
  MenuItem.prototype.unregister = function () {
    Spicetify.Menu._items = Spicetify.Menu._items.filter(function (x) { return x !== this; }, this);
    Spicetify.Menu._render();
  };
  MenuItem.prototype._build = function () {
    var self = this;
    return el("button", {
      className: "spotilify-menu-item",
      textContent: (this.props && this.props.name) || "",
      onclick: function (e) { try { self.props.onClick && self.props.onClick(e); } catch (err) { console.error(err); } }
    });
  };

  function MenuGroup(props) {
    if (!(this instanceof MenuGroup)) return new MenuGroup(props);
    this.props = props || {};
    this._items = (props.items || []).map(function (it) {
      return (it && it.register) ? it : MenuItem(it);
    });
  }
  MenuGroup.prototype.register = function () {
    var self = this;
    Spicetify.Menu._items.push(self);
    Spicetify.Menu._render();
    return function () { self.unregister(); };
  };
  MenuGroup.prototype.unregister = function () {
    Spicetify.Menu._items = Spicetify.Menu._items.filter(function (x) { return x !== this; }, this);
    Spicetify.Menu._render();
  };
  MenuGroup.prototype._build = function () {
    var wrap = el("div", { className: "spotilify-menu-group" });
    if (this.props.name) wrap.appendChild(el("div", { className: "spotilify-menu-group-title", textContent: this.props.name }));
    this._items.forEach(function (it) { wrap.appendChild(it._build()); });
    return wrap;
  };

  window.Spicetify.Menu = {
    Item: MenuItem,
    ItemSubMenu: MenuGroup,
    _items: [],
    _container: null,
    _render: function () {
      var c = Spicetify.Menu._container;
      if (!c) {
        c = el("div", { className: "spotilify-menu-root", style: { display: "none" } });
        document.body.appendChild(c);
        Spicetify.Menu._container = c;
      }
      c.innerHTML = "";
      Spicetify.Menu._items.forEach(function (it) { c.appendChild(it._build()); });
    },
    show: function (x, y) {
      var c = Spicetify.Menu._container;
      if (!c) return;
      c.style.left = (x || 0) + "px";
      c.style.top = (y || 0) + "px";
      c.style.display = "block";
    },
    hide: function () {
      var c = Spicetify.Menu._container;
      if (c) c.style.display = "none";
    }
  };
  document.addEventListener("click", function (e) {
    if (Spicetify.Menu && Spicetify.Menu._container && !Spicetify.Menu._container.contains(e.target)) {
      Spicetify.Menu.hide();
    }
  });

  // ContextMenu: same MenuItem shape; combos aggregate items shown together.
  var ContextMenuItem = MenuItem;

  function ContextMenuCombo(props) {
    if (!(this instanceof ContextMenuCombo)) return new ContextMenuCombo(props);
    this.props = props || {};
    this._items = (props.items || []).map(function (it) {
      return (it && it.register) ? it : ContextMenuItem(it);
    });
  }
  ContextMenuCombo.prototype.register = function () {
    var self = this;
    Spicetify.ContextMenu._combos.push(self);
    return function () {
      Spicetify.ContextMenu._combos = Spicetify.ContextMenu._combos.filter(function (x) { return x !== self; });
    };
  };

  window.Spicetify.ContextMenu = {
    Item: ContextMenuItem,
    Combo: ContextMenuCombo,
    _combos: [],
    _buildFor: function () {
      var wrap = el("div", { className: "spotilify-contextmenu" });
      Spicetify.ContextMenu._combos.forEach(function (combo) {
        combo._items.forEach(function (it) { wrap.appendChild(it._build()); });
      });
      return wrap;
    }
  };

  // Right-click interceptor: when user right-clicks a track-looking row,
  // surface our entries in a floating menu near the cursor.
  document.addEventListener("contextmenu", function (e) {
    if (!Spicetify.ContextMenu._combos.length) return;
    var row = e.target.closest && e.target.closest("[data-testid='tracklist-row'], [href*='/track/']");
    if (!row) return;
    e.preventDefault();
    var menu = Spicetify.ContextMenu._buildFor();
    menu.style.position = "fixed";
    menu.style.left = e.clientX + "px";
    menu.style.top = e.clientY + "px";
    menu.style.zIndex = "2147483646";
    menu.style.background = "#282828";
    menu.style.border = "1px solid rgba(255,255,255,0.1)";
    menu.style.borderRadius = "8px";
    menu.style.padding = "6px 0";
    menu.style.minWidth = "180px";
    document.body.appendChild(menu);
    var remove = function () { menu.remove(); document.removeEventListener("click", remove); };
    setTimeout(function () { document.addEventListener("click", remove); }, 0);
  }, true);

  // ---- Topbar / Playbar buttons -------------------------------------------

  function makeButton(opts) {
    return el("button", {
      className: "spotilify-bar-btn",
      title: opts.label || "",
      innerHTML: opts.icon || "",
      onclick: function (e) { try { opts.onClick && opts.onClick(e); } catch (err) { console.error(err); } }
    });
  }

  function TopbarButton(opts) {
    if (!(this instanceof TopbarButton)) return new TopbarButton(opts);
    this.opts = opts || {};
    this._node = null;
  }
  TopbarButton.prototype.register = function () {
    var self = this;
    function attach() {
      var container = findSelector(".main-topBar-container, .main-topbar-container");
      if (!container) { setTimeout(attach, 500); return; }
      var b = makeButton(self.opts);
      b.classList.add("spotilify-topbar-btn");
      container.appendChild(b);
      self._node = b;
      Spicetify.Topbar._buttons.push(self);
    }
    attach();
    return function () { self.unregister(); };
  };
  TopbarButton.prototype.unregister = function () {
    if (this._node && this._node.parentNode) this._node.parentNode.removeChild(this._node);
    Spicetify.Topbar._buttons = Spicetify.Topbar._buttons.filter(function (x) { return x !== this; }, this);
  };
  window.Spicetify.Topbar = { Button: TopbarButton, _buttons: [] };

  function PlaybarButton(opts) {
    if (!(this instanceof PlaybarButton)) return new PlaybarButton(opts);
    this.opts = opts || {};
    this._node = null;
  }
  PlaybarButton.prototype.register = function () {
    var self = this;
    function attach() {
      var container = findSelector(".main-nowPlayingBar-extraControls, .main-nowPlayingBar-rightButtons");
      if (!container) { setTimeout(attach, 500); return; }
      var b = makeButton(self.opts);
      b.classList.add("spotilify-playbar-btn");
      container.appendChild(b);
      self._node = b;
      Spicetify.Playbar._buttons.push(self);
    }
    attach();
    return function () { self.unregister(); };
  };
  PlaybarButton.prototype.unregister = function () {
    if (this._node && this._node.parentNode) this._node.parentNode.removeChild(this._node);
    Spicetify.Playbar._buttons = Spicetify.Playbar._buttons.filter(function (x) { return x !== this; }, this);
  };
  window.Spicetify.Playbar = { Button: PlaybarButton, _buttons: [] };

  // ---- PopupModal ---------------------------------------------------------
  window.Spicetify.PopupModal = {
    display: function (opts) {
      var overlay = el("div", { className: "spotilify-modal-overlay" });
      overlay.style.cssText = "position:fixed;inset:0;background:rgba(0,0,0,0.7);z-index:2147483647;display:flex;align-items:center;justify-content:center;";
      var dialog = el("div", { className: "spotilify-modal-dialog" });
      dialog.style.cssText = "position:relative;background:#282828;border-radius:12px;max-width:520px;width:90vw;max-height:80vh;overflow:auto;padding:24px;color:#fff;";
      if (opts.title) dialog.appendChild(el("h2", { textContent: opts.title, style: { margin: "0 0 16px 0" } }));
      if (opts.content) {
        if (typeof opts.content === "string") dialog.insertAdjacentHTML("beforeend", opts.content);
        else dialog.appendChild(opts.content);
      }
      var close = el("button", { textContent: "×", style: { position: "absolute", top: "8px", right: "12px", background: "transparent", color: "#fff", border: "none", fontSize: "22px", cursor: "pointer" } });
      close.onclick = function () { overlay.remove(); };
      dialog.appendChild(close);
      overlay.appendChild(dialog);
      overlay.addEventListener("click", function (e) { if (e.target === overlay) overlay.remove(); });
      document.body.appendChild(overlay);
      return overlay;
    },
    hide: function () {
      var existing = document.querySelector(".spotilify-modal-overlay");
      if (existing) existing.remove();
    }
  };

  // ---- Panel --------------------------------------------------------------
  window.Spicetify.Panel = {
    register: function (props) {
      var id = props.id || ("panel-" + Math.random().toString(36).slice(2));
      var container = el("div", { id: id, className: "spotilify-panel" });
      container.style.cssText = "display:none;position:fixed;right:0;top:32px;bottom:0;width:320px;background:#181818;z-index:9999;color:#fff;padding:16px;overflow:auto;";
      document.body.appendChild(container);
      if (typeof props.render === "function") {
        try { props.render(container); } catch (e) { console.error("[Spotilite] Panel render", e); }
      }
      return {
        show: function () { container.style.display = "block"; },
        hide: function () { container.style.display = "none"; },
        toggle: function () { container.style.display = container.style.display === "none" ? "block" : "none"; },
        remove: function () { container.remove(); }
      };
    }
  };

  // ---- ReactComponent stubs ----------------------------------------------
  function stub(componentName) {
    return function (props) {
      return el("div", { className: "spotilify-stub spotilify-stub--" + componentName, textContent: (props && props.children) || "" });
    };
  }
  window.Spicetify.ReactComponent = {
    ButtonPrimary: stub("ButtonPrimary"),
    ButtonSecondary: stub("ButtonSecondary"),
    ButtonTertiary: stub("ButtonTertiary"),
    Menu: stub("Menu"),
    MenuItem: stub("MenuItem"),
    TooltipWrapper: stub("TooltipWrapper"),
    TextComponent: stub("TextComponent"),
    IconComponent: stub("IconComponent"),
    ConfirmDialog: stub("ConfirmDialog"),
    Slider: stub("Slider"),
    Toggle: stub("Toggle"),
    Dropdown: stub("Dropdown"),
    Router: stub("Router"),
    Routes: stub("Routes"),
    Route: stub("Route")
  };

  // Minimal React shim so extensions that just call Spicetify.React.createElement
  // defensively won't throw. This is NOT a real React: it constructs DOM nodes
  // directly. Hooks/concurrent features are not supported.
  window.Spicetify.React = {
    createElement: function (type, props) {
      var children = Array.prototype.slice.call(arguments, 2);
      return el(typeof type === "string" ? type : "div", props, children);
    },
    Fragment: "Fragment",
    Component: function () {},
    isValidElement: function () { return false; }
  };
  window.Spicetify.ReactDOM = {
    render: function (node, container) {
      if (container) {
        container.innerHTML = "";
        if (node) container.appendChild(node);
      }
    },
    unmountComponentAtNode: function (container) { if (container) container.innerHTML = ""; }
  };

  // Style block for our injected UI primitives.
  function injectStyles() {
    if (document.getElementById("spotilify-component-styles")) return;
    var s = document.createElement("style");
    s.id = "spotilify-component-styles";
    s.textContent = [
      ".spotilify-menu-root { position:fixed; background:#282828; border:1px solid rgba(255,255,255,0.1); border-radius:8px; padding:6px 0; min-width:200px; z-index:2147483646; box-shadow:0 8px 24px rgba(0,0,0,0.6); }",
      ".spotilify-menu-item, .spotilify-contextmenu button { display:block; width:100%; text-align:left; padding:8px 14px; background:transparent; border:none; color:#ddd; font-size:13px; cursor:pointer; font-family:inherit; }",
      ".spotilify-menu-item:hover, .spotilify-contextmenu button:hover { background:rgba(255,255,255,0.08); color:#fff; }",
      ".spotilify-menu-group-title { padding:6px 14px 4px; color:#1db954; font-size:10px; font-weight:700; text-transform:uppercase; }",
      ".spotilify-bar-btn { background:transparent; border:none; color:#b3b3b3; cursor:pointer; padding:6px 8px; border-radius:4px; }",
      ".spotilify-bar-btn:hover { color:#fff; background:rgba(255,255,255,0.08); }",
      ".spotilify-topbar-btn, .spotilify-playbar-btn { margin-left:8px; }"
    ].join("\n");
    document.head.appendChild(s);
  }
  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", injectStyles);
  else injectStyles();
})();
