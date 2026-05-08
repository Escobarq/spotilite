package spotify

const titleBarScript = `(function() {
if (window.__spotilite_loaded) return;
window.__spotilite_loaded = true;
if (!window.__origFetch) window.__origFetch = window.fetch.bind(window);
var API = 'http://localhost:8765';
var lang = localStorage.getItem('spotilite.lang') || 'es';
var adBlock = localStorage.getItem('spotilite.adblock') !== 'false';
var sectionBlock = localStorage.getItem('spotilite.sectionblock') !== 'false';
var premiumSpoof = localStorage.getItem('spotilite.premium_spoof') !== 'false';
var experiments = localStorage.getItem('spotilite.experiments') !== 'false';
var historyOn = localStorage.getItem('spotilite.history') !== 'false';
var bgMode = localStorage.getItem('spotilite.bg') !== 'false';
var txt = lang === 'en'
  ? {ad:'Ad Blocker',sec:'Block Sections',prem:'Hide Premium',exp:'Experiments',hist:'History',set:'Settings',lang:'Language',es:'Spanish',en:'English',app:'App',bg:'Background',min:'Minimize',max:'Maximize',rest:'Restore',close:'Close'}
  : {ad:'Bloquear Anuncios',sec:'Bloquear Secciones',prem:'Ocultar Premium',exp:'Experimentos',hist:'Historial',set:'Ajustes',lang:'Idioma',es:'Espanol',en:'English',app:'App',bg:'Segundo Plano',min:'Minimizar',max:'Maximizar',rest:'Restaurar',close:'Cerrar'};
function apiPost(path, data) {
  window.__origFetch(API + path, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(data || {})
  }).catch(function(e){});
}
var bar = document.getElementById('spotilite-title-bar');
if (bar) return;
bar = document.createElement('div');
bar.id = 'spotilite-title-bar';
bar.style.cssText = 'display:flex;align-items:center;';
document.body.insertBefore(bar, document.body.firstChild);
var left = document.createElement('div');
left.className = 'spotilite-left';
left.innerHTML = '<span class="spotilite-logo">S</span><span class="spotilite-title">Spotilite</span>';
bar.appendChild(left);
var center = document.createElement('div');
center.className = 'spotilite-center';
center.style.cssText = 'flex:1;display:flex;align-items:center;justify-content:center;height:100%;--wails-draggable:drag;';
var toggles = document.createElement('div');
toggles.className = 'spotilite-toggles';
function mkBtn(id, active, title, svg) {
  var b = document.createElement('button');
  b.id = id;
  b.className = 'spotilite-icon-btn' + (active ? ' active' : '');
  b.title = title;
  b.innerHTML = svg;
  return b;
}
var adBtn = mkBtn('spotilite-toggle-adblock', adBlock, txt.ad, '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>');
var secBtn = mkBtn('spotilite-toggle-sectionblock', sectionBlock, txt.sec, '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M3 3h18v2H3V3zm0 16h18v2H3v-2zm0-8h18v2H3v-2z"/></svg>');
var premBtn = mkBtn('spotilite-toggle-premium', premiumSpoof, txt.prem, '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z"/></svg>');
var expBtn = mkBtn('spotilite-toggle-experiments', experiments, txt.exp, '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M7 2v11h3v9l7-12h-4l4-8z"/></svg>');
var histBtn = mkBtn('spotilite-toggle-history', historyOn, txt.hist, '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M13 3a9 9 0 00-9 9H1l3.89 3.89.07.14L9 12H6c0-3.87 3.13-7 7-7s7 3.13 7 7-3.13 7-7 7a7 7 0 01-5.19-2.32l-1.41 1.41A9 9 0 0013 21a9 9 0 000-18zm-1 5v5l4.28 2.54.72-1.21-3.5-2.08V8H12z"/></svg>');
toggles.appendChild(adBtn);
toggles.appendChild(secBtn);
toggles.appendChild(premBtn);
toggles.appendChild(expBtn);
toggles.appendChild(histBtn);
center.appendChild(toggles);
bar.appendChild(center);
var right = document.createElement('div');
right.className = 'spotilite-right';
var setCont = document.createElement('div');
setCont.style.cssText = 'position:relative;height:100%;';
var setBtn = document.createElement('button');
setBtn.className = 'spotilite-icon-btn spotilite-settings-btn';
setBtn.title = txt.set;
setBtn.innerHTML = '<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M19.14 12.94c.04-.31.06-.63.06-.94 0-.31-.02-.63-.06-.94l2.03-1.58a.49.49 0 00.12-.61l-1.92-3.32a.49.49 0 00-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.48.48 0 00-.48-.41h-3.84a.48.48 0 00-.48.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96a.49.49 0 00-.59.22L2.74 8.87a.49.49 0 00.12.61l2.03 1.58c-.05.31-.07.63-.07.94s.02.64.07.94l-2.03 1.58a.49.49 0 00-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.48-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32a.49.49 0 00-.12-.61l-2.01-1.58zM12 15.6A3.6 3.6 0 1115.6 12 3.6 3.6 0 0112 15.6z"/></svg>';
setCont.appendChild(setBtn);
var dd = document.createElement('div');
dd.className = 'spotilite-dropdown';
dd.innerHTML =
  '<div class="spotilite-dropdown-header">' + txt.lang + '</div>' +
  '<button class="spotilite-dropdown-item" data-lang="es"><span class="spotilite-lang-label">ES</span><span>' + txt.es + '</span><span class="spotilite-radio' + (lang==='es'?' on':'') + '"></span></button>' +
  '<button class="spotilite-dropdown-item" data-lang="en"><span class="spotilite-lang-label">EN</span><span>' + txt.en + '</span><span class="spotilite-radio' + (lang==='en'?' on':'') + '"></span></button>' +
  '<div class="spotilite-dropdown-sep"></div>' +
  '<div class="spotilite-dropdown-header">' + txt.app + '</div>' +
  '<button class="spotilite-dropdown-item" id="spotilite-bg-toggle"><span>' + txt.bg + '</span><span class="spotilite-toggle' + (bgMode?' on':'') + '"></span></button>';
setCont.appendChild(dd);
right.appendChild(setCont);
var winBtns = document.createElement('div');
winBtns.className = 'spotilite-win-btns';
winBtns.innerHTML = '<button class="spotilite-win-btn" id="spotilite-minimize">&#8722;</button><button class="spotilite-win-btn" id="spotilite-maximize">&#9633;</button><button class="spotilite-win-btn spotilite-close" id="spotilite-close">&#10005;</button>';
right.appendChild(winBtns);
bar.appendChild(right);
setBtn.onclick = function(e) {
  e.stopPropagation();
  dd.classList.toggle('active');
};
document.addEventListener('click', function(e) {
  if (!setCont.contains(e.target)) dd.classList.remove('active');
});
dd.querySelector('[data-lang="es"]').onclick = function() {
  lang='es'; localStorage.setItem('spotilite.lang','es'); apiPost('/api/settings/lang',{lang:'es'}); location.reload();
};
dd.querySelector('[data-lang="en"]').onclick = function() {
  lang='en'; localStorage.setItem('spotilite.lang','en'); apiPost('/api/settings/lang',{lang:'en'}); location.reload();
};
document.getElementById('spotilite-bg-toggle').onclick = function() {
  bgMode=!bgMode; localStorage.setItem('spotilite.bg',bgMode); apiPost('/api/settings/background',{enabled:bgMode}); upd();
};
adBtn.onclick = function() { adBlock=!adBlock; localStorage.setItem('spotilite.adblock',adBlock); apiPost('/api/spotx/module',{module:'adblock',enabled:adBlock}); upd(); };
secBtn.onclick = function() { sectionBlock=!sectionBlock; localStorage.setItem('spotilite.sectionblock',sectionBlock); apiPost('/api/spotx/module',{module:'sectionblock',enabled:sectionBlock}); upd(); };
premBtn.onclick = function() { premiumSpoof=!premiumSpoof; localStorage.setItem('spotilite.premium_spoof',premiumSpoof); apiPost('/api/spotx/module',{module:'premium_spoof',enabled:premiumSpoof}); upd(); };
expBtn.onclick = function() { experiments=!experiments; localStorage.setItem('spotilite.experiments',experiments); apiPost('/api/spotx/module',{module:'experiments',enabled:experiments}); upd(); };
histBtn.onclick = function() { historyOn=!historyOn; localStorage.setItem('spotilite.history',historyOn); apiPost('/api/spotx/module',{module:'history',enabled:historyOn}); upd(); };
function upd() {
  adBtn.className='spotilite-icon-btn'+(adBlock?' active':'');
  secBtn.className='spotilite-icon-btn'+(sectionBlock?' active':'');
  premBtn.className='spotilite-icon-btn'+(premiumSpoof?' active':'');
  expBtn.className='spotilite-icon-btn'+(experiments?' active':'');
  histBtn.className='spotilite-icon-btn'+(historyOn?' active':'');
}
upd();
document.getElementById('spotilite-minimize').onclick = function() { apiPost('/api/window/minimize', {}); };
document.getElementById('spotilite-maximize').onclick = function() { apiPost('/api/window/maximize', {}); };
document.getElementById('spotilite-close').onclick = function() { apiPost('/api/window/close', {}); };
function applyOffset() {
  var barH = 28;
  var main = document.querySelector('.Root') || document.querySelector('.main-view-container') || document.querySelector('#main');
  if (main) main.style.paddingTop = barH + 'px';
  var topBar = document.querySelector('.main-topBar-container');
  if (topBar) topBar.style.display = 'none';
  var resp = document.querySelector('[class*="response"]') || document.querySelector('[class*="mobile"]') || document.querySelector('[class*="touch"]');
  if (resp) resp.style.paddingTop = barH + 'px';
}
applyOffset();
setInterval(applyOffset, 2000);
window.addEventListener('resize', function() { applyOffset(); });
})();`

const titleBarCSS = `
html, body { margin: 0 !important; padding: 0 !important; overflow: hidden !important; }
#spotilite-title-bar {
  position: fixed; top: 0; left: 0; width: 100%; height: 28px;
  z-index: 2147483647; background: #191414;
  display: flex; align-items: center; justify-content: space-between;
  --wails-draggable: drag;
  user-select: none;
}
#spotilite-title-bar .spotilite-left {
  display: flex; align-items: center; gap: 8px;
  padding-left: 12px; height: 100%;
  --wails-draggable: drag;
  flex-shrink: 0;
}
#spotilite-title-bar .spotilite-center {
  flex: 1; display: flex; align-items: center; justify-content: center;
  height: 100%; --wails-draggable: drag; cursor: grab;
}
#spotilite-title-bar .spotilite-center:active { cursor: grabbing; }
@media screen and (max-width: 768px) {
  #spotilite-title-bar .spotilite-toggles { gap: 0; }
  #spotilite-title-bar .spotilite-icon-btn { padding: 0 4px; }
  #spotilite-title-bar .spotilite-title { display: none; }
}
#spotilite-title-bar .spotilite-toggles {
  display: flex; align-items: center; height: 100%; gap: 2px;
}
#spotilite-title-bar .spotilite-logo {
  width: 18px; height: 18px; background: #1DB954;
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: bold; color: #000;
}
#spotilite-title-bar .spotilite-title {
  color: #fff; font-size: 12px; font-weight: 600;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}
#spotilite-title-bar .spotilite-toggles {
  display: flex; align-items: center; height: 100%; gap: 2px;
}
#spotilite-title-bar .spotilite-icon-btn {
  height: 100%; padding: 0 6px; background: transparent; border: none;
  color: #666; font-size: 13px; cursor: pointer;
  transition: color 0.15s, background 0.15s;
  display: flex; align-items: center; justify-content: center;
  pointer-events: auto;
}
#spotilite-title-bar .spotilite-icon-btn:hover { color: #fff; background: rgba(255,255,255,0.08); }
#spotilite-title-bar .spotilite-icon-btn.active { color: #1DB954; }
#spotilite-title-bar .spotilite-icon-btn.active:hover { color: #1ed760; }
#spotilite-title-bar .spotilite-settings-btn { color: #888; }
#spotilite-title-bar .spotilite-settings-btn:hover { color: #fff; }
#spotilite-title-bar .spotilite-right {
  display: flex; align-items: center; height: 100%;
}
#spotilite-title-bar .spotilite-win-btns {
  display: flex; height: 100%;
  pointer-events: auto;
}
#spotilite-title-bar .spotilite-win-btn {
  width: 46px; height: 100%; background: transparent; border: none;
  color: #b3b3b3; font-size: 14px; cursor: pointer;
  transition: background 0.15s, color 0.15s;
  pointer-events: auto;
}
#spotilite-title-bar .spotilite-dropdown {
  position: absolute; top: 100%; right: 0;
  background: #282828; border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px; box-shadow: 0 8px 24px rgba(0,0,0,0.6);
  z-index: 2147483647; width: 200px; padding: 6px 0;
  display: none; max-height: 80vh; overflow-y: auto;
  pointer-events: auto;
}
#spotilite-title-bar .spotilite-dropdown.active { display: block; }
#spotilite-title-bar .spotilite-dropdown::-webkit-scrollbar { width: 6px; }
#spotilite-title-bar .spotilite-dropdown::-webkit-scrollbar-thumb { background: #555; border-radius: 3px; }
#spotilite-title-bar .spotilite-dropdown-header {
  padding: 8px 14px 4px; color: #1DB954; font-size: 10px; font-weight: 700;
  text-transform: uppercase; letter-spacing: 0.5px;
}
#spotilite-title-bar .spotilite-dropdown-sep {
  height: 1px; background: rgba(255,255,255,0.08); margin: 6px 0;
}
#spotilite-title-bar .spotilite-dropdown-item {
  width: 100%; text-align: left; padding: 7px 14px; background: transparent;
  border: none; color: #ddd; font-size: 12px; cursor: pointer;
  font-family: inherit; display: flex; align-items: center; justify-content: space-between;
  transition: background 0.1s;
}
#spotilite-title-bar .spotilite-dropdown-item:hover { background: rgba(255,255,255,0.08); color: #fff; }
#spotilite-title-bar .spotilite-lang-label {
  display: inline-flex; align-items: center; justify-content: center;
  width: 22px; height: 14px; border-radius: 2px;
  background: #444; color: #fff; font-size: 9px; font-weight: 700; margin-right: 6px;
}
#spotilite-title-bar .spotilite-toggle {
  width: 28px; height: 16px; border-radius: 8px; background: #555;
  position: relative; transition: background 0.2s; flex-shrink: 0;
}
#spotilite-title-bar .spotilite-toggle::after {
  content: ''; position: absolute; top: 2px; left: 2px;
  width: 12px; height: 12px; border-radius: 50%;
  background: #fff; transition: transform 0.2s;
}
#spotilite-title-bar .spotilite-toggle.on { background: #1DB954; }
#spotilite-title-bar .spotilite-toggle.on::after { transform: translateX(12px); }
#spotilite-title-bar .spotilite-radio {
  width: 12px; height: 12px; border-radius: 50%; border: 2px solid #555;
  position: relative; flex-shrink: 0; transition: border-color 0.2s;
}
#spotilite-title-bar .spotilite-radio.on { border-color: #1DB954; }
#spotilite-title-bar .spotilite-radio.on::after {
  content: ''; position: absolute; top: 2px; left: 2px;
  width: 4px; height: 4px; border-radius: 50%; background: #1DB954;
}
`