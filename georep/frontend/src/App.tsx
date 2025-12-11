import { useEffect, useState } from "react";
import ReaderPage from "./pages/ReaderPage";
import WriterPage from "./pages/WriterPage";
import { Session } from "./types";

export default function App() {
  const [autoRegion, setAutoRegion] = useState("eu");
  const [manualRegion, setManualRegion] = useState<string | null>(null);
  const [mode, setMode] = useState<"select" | "reader" | "writer">("select");
  const [selectedRegion, setSelectedRegion] = useState("eu");
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [showPasswordInput, setShowPasswordInput] = useState(false);
  const [writerPassword, setWriterPassword] = useState("");
  const [error, setError] = useState("");

  const API_BASE = "http://localhost:8080/api";

  useEffect(() => {
    fetch(`${API_BASE}/region`)
      .then((r) => r.json())
      .then((d) => {
        if (d.region && d.region !== "unknown") setAutoRegion(d.region);
      });
  }, []);

  const handleArticleAdded = () => setRefreshTrigger((p) => p + 1);

  const handleWriterClick = () => {
    setShowPasswordInput(true);
    setError("");
  };

  const handleWriterConfirm = () => {
    if (writerPassword === "4321") {
      setManualRegion(selectedRegion);
      setMode("writer");
    } else {
      setError("Yanlış parola! Lütfen tekrar deneyin.");
    }
  };

  if (mode === "select") {
    const regions = ["eu", "us", "asia", "sa", "africa"];
    return (
      <div style={centerBox}>
        <h1>🌍 Geo-Replicated News Portal</h1>
        <p>
          Otomatik tespit edilen bölgeniz:{" "}
          <b style={{ color: "#3182ce" }}>{autoRegion.toUpperCase()}</b>
        </p>

        <div style={regionBox}>
          <h3>Test İçin Bölge Seçimi</h3>
          <select
            value={selectedRegion}
            onChange={(e) => setSelectedRegion(e.target.value)}
            style={selectStyle}
          >
            {regions.map((r) => (
              <option key={r} value={r}>
                {r.toUpperCase()}
              </option>
            ))}
          </select>
        </div>

        <div style={{ marginTop: "20px" }}>
          <button
            style={button("#22c55e")}
            onClick={() => {
              setManualRegion(selectedRegion);
              setMode("reader");
            }}
          >
            👀 Okuyucu
          </button>

          <button style={button("#3b82f6")} onClick={handleWriterClick}>
            ✍️ Yazar
          </button>
        </div>

        {/* 🔐 Parola alanı sadece Yazar tıklanınca görünür */}
        {showPasswordInput && (
          <div style={passwordBox}>
            <input
              type="password"
              placeholder="Yazar parolası (örnek: 1234)"
              value={writerPassword}
              onChange={(e) => setWriterPassword(e.target.value)}
              style={passwordInput}
            />
            <button onClick={handleWriterConfirm} style={button("#3b82f6")}>
              Giriş
            </button>
          </div>
        )}

        {error && <p style={{ color: "red", marginTop: "8px" }}>{error}</p>}
      </div>
    );
  }

  const currentRegion = manualRegion || autoRegion;
  const session: Session = {
    username: "demo",
    role: mode === "writer" ? "writer" : "reader",
    region: currentRegion as any,
    token: "demo-token",
  };

  if (mode === "reader")
    return (
      <ReaderPage
        session={session}
        onLogout={() => setMode("select")}
        refreshTrigger={refreshTrigger}
      />
    );

  if (mode === "writer")
    return (
      <WriterPage
        session={session}
        onLogout={() => setMode("select")}
        onArticleAdded={handleArticleAdded}
      />
    );

  return null;
}

/* Styles */
const centerBox = {
  display: "flex",
  flexDirection: "column" as const,
  justifyContent: "center",
  alignItems: "center",
  height: "100vh",
  backgroundColor: "#f5f7fa",
};
const regionBox = {
  background: "white",
  padding: 20,
  borderRadius: 10,
  marginTop: 20,
};
const selectStyle = {
  padding: 8,
  borderRadius: 6,
  border: "1px solid #ccc",
};
const button = (bg: string) => ({
  backgroundColor: bg,
  color: "white",
  border: "none",
  padding: "10px 20px",
  borderRadius: "8px",
  margin: "10px",
  cursor: "pointer",
});
const passwordBox = {
  marginTop: "15px",
  display: "flex",
  flexDirection: "column" as const,
  alignItems: "center",
};
const passwordInput = {
  padding: "8px",
  marginBottom: "10px",
  borderRadius: "6px",
  border: "1px solid #ccc",
  width: "200px",
};
