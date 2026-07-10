package spicetify

import (
	"encoding/json"
	"fmt"
	"strings"

	"spotilite/internal/customapps"
)

// WebpackChunkTemplate is the shape used by spicetify-cli to register a custom
// app as a webpack chunk on the client. See cli-main/src/cmd/apply.go:369-379.
//
// Desktop push:
//
//	((globalScope)=>((globalScope.rspackChunkclient_web =
//	  globalScope.rspackChunkclient_web ||
//	  globalScope.webpackChunkclient_web || []),
//	  (globalScope.webpackChunkclient_web =
//	    globalScope.webpackChunkclient_web ||
//	    globalScope.rspackChunkclient_web)))
//	("undefined"!=typeof self?self:global)
//	.push([["<appName>"],{"<appName>":(e,t,n)=>{
//	  "use strict";n.r(t),n.d(t,{default:()=>render});
//	  <index.js body>
//	}}]);
//
// On the Spotify web player `webpackChunkclient_web` already exists on window,
// so this push registers a new named entry that the router can lazy-load for
// route "/<appName>". The companion RouteHook does the mount work the
// desktop client router patch would have done.
func WebpackChunkTemplate(app *customapps.App, indexJSBody string) string {
	if app == nil {
		return ""
	}
	folder := app.FolderName
	if folder == "" {
		folder = app.Name
	}
	body := indexJSBody
	if body == "" {
		body = app.IndexJS
	}
	return fmt.Sprintf(webpackPushSkeleton, encodeStringLiteral(folder), encodeStringLiteral(folder), body)
}

const webpackPushSkeleton = `((globalScope)=>((globalScope.rspackChunkclient_web=globalScope.rspackChunkclient_web||globalScope.webpackChunkclient_web||[]),(globalScope.webpackChunkclient_web=globalScope.webpackChunkclient_web||globalScope.rspackChunkclient_web)))("undefined"!=typeof self?self:global).push([[%s],{%s:(e,t,n)=>{"use strict";n.r(t),n.d(t,{default:()=>render});
%s
}}]);`

// RouteHook generates the JS that hijacks the Spotify router to mount a custom
// app at /<folderName>. The desktop patch (cli-main preprocess.go) wraps the
// router navigate function to intercept custom app routes; we cannot patch the
// bundle here, so we listen on history.pushState/popstate and, when the route
// matches "/<folderName>", inject the app's render(container) into a fresh
// DOM node. We also surface the app's icon in the sidebar (mirrors the desktop
// sidebar injection).
func RouteHook(app *customapps.App) string {
	if app == nil || app.Manifest == nil {
		return ""
	}
	folder := app.FolderName
	if folder == "" {
		folder = app.Name
	}
	display := app.Name
	icon := app.Manifest.Icon
	activeIcon := app.Manifest.ActiveIcon
	if activeIcon == "" {
		activeIcon = icon
	}
	displayJSONBytes, _ := json.Marshal(display)
	iconJSONBytes, _ := json.Marshal(icon)
	folderJSONBytes, _ := json.Marshal(folder)
	displayJSON := string(displayJSONBytes)
	iconJSON := string(iconJSONBytes)
	folderJSON := string(folderJSONBytes)
	cssInjection := ""
	if app.StylesCSS != "" {
		cssInjection = "var st=document.createElement('style');st.id='spotilite-app-" + folder + "-css';st.textContent=" + jsStringSp(app.StylesCSS) + ";document.head.appendChild(st);"
	}
	var subfiles strings.Builder
	for _, s := range app.Subfiles {
		subfiles.WriteString(s)
		subfiles.WriteString("\n")
	}
	subfilesJSONBytes, _ := json.Marshal(subfiles.String())
	subfilesJSON := string(subfilesJSONBytes)

	var b strings.Builder
	b.WriteString("/* === spotilite custom app: ")
	b.WriteString(folder)
	b.WriteString(" === */\n(function(){\n")
	b.WriteString("var folder=" + folderJSON + ";\n")
	b.WriteString("var display=" + displayJSON + ";\n")
	b.WriteString("var icon=" + iconJSON + ";\n")
	b.WriteString("var route='/' + folder;\n")
	b.WriteString("var containerId='spotilite-app-' + folder;\n")
	b.WriteString(cssInjection)
	b.WriteString("\n")
	b.WriteString("function logM(m){console.log('[Spotilite CustomApp ' + folder + '] ' + m);}\n")
	b.WriteString("function makeContainer(){var c=document.getElementById(containerId);if(c)return c;c=document.createElement('div');c.id=containerId;c.style.cssText='position:absolute;top:0;left:0;right:0;bottom:0;z-index:9999;background:var(--background-base,#000)';document.body.appendChild(c);return c;}\n")
	b.WriteString("var injected=false;\n")
	b.WriteString("function injectSidebarEntry(){if(injected)return;var sel=['li.main-navBar-navBarItem','.main-navBar-navBarItem','[class*=main-navBar-navBarItem]','.main-navBar-navBar','.main-appShell-sideBar li'];var n=null;for(var i=0;i<sel.length;i++){var e=document.querySelector(sel[i]);if(e){n=e;break;}}if(!n){logM('sidebar not ready, retrying');setTimeout(injectSidebarEntry,1500);return;}var p=n.parentElement||n;var li=document.createElement('li');li.className='main-navBar-navBarItem spotilite-app-link';li.dataset.name=folder;li.title=display;li.style.cssText='cursor:pointer;display:flex;align-items:center;justify-content:center;padding:12px;color:var(--spice-text,#fff)';li.innerHTML=icon;li.onclick=function(e){e.preventDefault();e.stopPropagation();try{window.history.pushState({spotiliteApp:folder},'',route);var c=makeContainer();mount(c);}catch(err){console.error('[Spotilite CustomApp '+folder+'] mount error:',err);}};p.appendChild(li);injected=true;logM('sidebar entry injected');}\n")
	b.WriteString("var subfilesJS=" + string(subfilesJSON) + ";\n")
	b.WriteString("function mount(container){logM('mounting');if(subfilesJS){try{eval(subfilesJS);}catch(err){console.error('[Spotilite CustomApp '+folder+'] subfiles error:',err);}}var renderFn=window.render||(function(){return null;});try{renderFn(container);}catch(err){console.error('[Spotilite CustomApp '+folder+'] render crashed:',err);}logM('mounted');}\n")
	b.WriteString("window.addEventListener('popstate',function(e){var st=e&&e.state||{};if(st.spotiliteApp==folder){var c=makeContainer();mount(c);}else{var c=document.getElementById(containerId);if(c)c.remove();}});\n")
	b.WriteString("setTimeout(injectSidebarEntry,2000);\n")
	b.WriteString("})();\n")
	return b.String()
}

// encodeStringLiteral returns folder name wrapped in JSON-style string quotes
// so it can be used in the webpack push array literal and object key. We pick
// JSON because webpackChunk.push accepts arbitrary chunk id shapes (strings
// or arrays of strings) and JSON gives us correct escaping.
func encodeStringLiteral(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// jsStringSp escapes a Go string for inclusion as a JS double-quoted string
// literal. Mirrors internal/spotify/injector.go:jsString so we don't depend
// on the spotify package from here.
func jsStringSp(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteByte(s[i])
		}
	}
	b.WriteByte('"')
	return b.String()
}
