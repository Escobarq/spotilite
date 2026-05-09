import React, { useState, useEffect } from "react";
import { IsBackgroundMode, SetBackgroundMode } from "../wailsjs/go/app/App";

function App() {
  const [backgroundMode, setBackgroundMode] = useState(false);
  const [showSettings, setShowSettings] = useState(false);

  useEffect(() => {
    IsBackgroundMode().then(setBackgroundMode);
  }, []);

  const handleToggleBackground = async () => {
    const newValue = !backgroundMode;
    await SetBackgroundMode(newValue);
    setBackgroundMode(newValue);
  };

  return (
    <div
      id="App"
      style={{
        width: "100%",
        height: "100vh",
        background: "#000",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        color: "#fff",
        fontFamily:
          '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
      }}
    >
      <button
        onClick={() => setShowSettings(!showSettings)}
        style={{
          position: "absolute",
          top: 16,
          right: 16,
          background: "none",
          border: "none",
          color: "#fff",
          fontSize: 24,
          cursor: "pointer",
        }}
        title="Settings"
      >
        &#9881;
      </button>

      <div>Loading Spotify...</div>

      {showSettings && (
        <div
          style={{
            position: "absolute",
            top: 60,
            right: 16,
            background: "#1e1e1e",
            borderRadius: 8,
            padding: 16,
            minWidth: 250,
            boxShadow: "0 4px 12px rgba(0,0,0,0.5)",
          }}
        >
          <h3 style={{ margin: "0 0 12px", fontSize: 16 }}>Settings</h3>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
            }}
          >
            <label style={{ fontSize: 14 }}>
              Ejecutar en segundo plano
            </label>
            <button
              onClick={handleToggleBackground}
              style={{
                width: 48,
                height: 26,
                borderRadius: 13,
                border: "none",
                background: backgroundMode ? "#1db954" : "#555",
                cursor: "pointer",
                position: "relative",
                transition: "background 0.2s",
              }}
            >
              <div
                style={{
                  width: 20,
                  height: 20,
                  borderRadius: "50%",
                  background: "#fff",
                  position: "absolute",
                  top: 3,
                  left: backgroundMode ? 25 : 3,
                  transition: "left 0.2s",
                }}
              />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
