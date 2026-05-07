package spotify

const titleBarScript = `(function() {
if (window.__spotilite_loaded) return;
window.__spotilite_loaded = true;
if (!window.__origFetch) window.__origFetch = window.fetch.bind(window);
var API = 'http://localhost:8765';
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
document.body.insertBefore(bar, document.body.firstChild);
var left = document.createElement('div');
left.className = 'spotilite-left';
left.innerHTML = '<span class="spotilite-logo">S</span><span class="spotilite-title">Spotilite</span>';
bar.appendChild(left);
var right = document.createElement('div');
right.className = 'spotilite-right';
right.innerHTML = '<button class="spotilite-btn" id="spotilite-minimize">&#8722;</button><button class="spotilite-btn" id="spotilite-maximize">&#9633;</button><button class="spotilite-btn spotilite-close" id="spotilite-close">&#10005;</button>';
right.style.cssText = 'display:flex;gap:0;height:100%;';
bar.appendChild(right);
document.getElementById('spotilite-minimize').onclick = function() { apiPost('/api/window/minimize', {}); };
document.getElementById('spotilite-maximize').onclick = function() { apiPost('/api/window/maximize', {}); };
document.getElementById('spotilite-close').onclick = function() { apiPost('/api/window/close', {}); };
function applyOffset() {
  var barH = 28;
  var main = document.querySelector('.Root') || document.querySelector('.main-view-container') || document.querySelector('#main');
  if (main) main.style.paddingTop = barH + 'px';
  var topBar = document.querySelector('.main-topBar-container');
  if (topBar) topBar.style.display = 'none';
}
applyOffset();
setInterval(applyOffset, 3000);
})();`