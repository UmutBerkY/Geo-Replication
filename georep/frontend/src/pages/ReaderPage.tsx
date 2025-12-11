import React, { useEffect, useState } from "react";
import { apiGet } from "../api";
import { Article, ReplicationStatus, Session } from "../types";

type Props = {
  session: Session;
  onLogout: () => void;
};

export default function ReaderPage({ session, onLogout }: Props) {
  const [articles, setArticles] = useState<Article[]>([]);
  const [status, setStatus] = useState<ReplicationStatus[]>([]);
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);

  // ⚡ GECİKME VERİSİ
  const [latencyText, setLatencyText] = useState<string | null>(null);
  const [loadingLatency, setLoadingLatency] = useState(false);

  const loadArticles = async () => {
    try {
      const data = await apiGet<Article[]>(`/articles?region=${session.region}`);
      setArticles(data || []);
    } catch {
      setArticles([]);
    }
  };

  const loadStatus = async () => {
    try {
      const s = await apiGet<ReplicationStatus[]>("/replication-status");
      setStatus((s || []).filter((r) => r.replica.toLowerCase() !== "eu"));
    } catch {
      setStatus([]);
    }
  };

  const loadLatency = async () => {
    setLoadingLatency(true);
    try {
      const res = await apiGet<{
        region: string;
        latency: string;
        measured: string;
      }>(`/latency?region=${session.region}`);
      setLatencyText(res.latency);
    } catch {
      setLatencyText("Gecikme ölçümü yapılamadı.");
    } finally {
      setLoadingLatency(false);
    }
  };

  useEffect(() => {
    loadArticles();
    loadStatus();
    loadLatency();

    const interval = setInterval(() => {
      loadArticles();
      loadStatus();
    }, 5000);

    return () => clearInterval(interval);
  }, [session.region]);

  return (
    <section className="card">
      <div className="header-row">
        <div>
          <h2>
            Okuyucu görünümü – Bölge: {session.region.toUpperCase()}{" "}
            {session.region.toLowerCase() === "eu"
              ? "(Master sunucu)"
              : "(en yakın replika)"}
          </h2>
          <p className="hint">
            Okuma istekleri seçtiğin bölgeye en yakın replika veya EU master’dan
            geliyor. Yazılar birkaç saniye gecikmeli güncellenebilir.
          </p>

          {latencyText && (
            <p
              className="hint"
              style={{
                marginTop: "0.5rem",
                fontWeight: 500,
                color: "#1e40af",
              }}
            >
              {loadingLatency ? "Gecikme ölçülüyor..." : latencyText}
            </p>
          )}
        </div>
        <button onClick={onLogout}>Çıkış</button>
      </div>

      <div className="actions" style={{ gap: "0.5rem" }}>
        <button onClick={loadArticles}>Yazıları yenile</button>
        <button onClick={loadStatus}>Replikasyon durumu</button>
        <button onClick={loadLatency} disabled={loadingLatency}>
          {loadingLatency ? "Ölçülüyor..." : "Gecikme ölç"}
        </button>
      </div>

      <div className="stories">
        {articles.length > 0 ? (
          articles.map((a) => (
            <div
              key={a.id}
              style={{
                border: "1px solid #e2e8f0",
                borderRadius: "8px",
                padding: "15px",
                marginBottom: "10px",
                cursor: "pointer",
                background: "#fff",
                transition: "0.2s",
              }}
              onClick={() => setSelectedArticle(a)}
            >
              <h3 style={{ marginBottom: "6px" }}>{a.title}</h3>
              <p style={{ fontSize: "14px", color: "#4a5568" }}>
                {a.summary || a.content}
              </p>
              <p style={{ fontSize: "13px", color: "#718096", marginTop: "6px" }}>
                ✍️ {a.author}
              </p>
            </div>
          ))
        ) : (
          <p className="hint">
            Bu replikada henüz veri yok. Bir yazar EU master üzerinden yeni
            haber yayınladığında burada görünecek.
          </p>
        )}
      </div>

      {/* 📖 MAKALENİN UZUN DETAY GÖRÜNÜMÜ */}
      {selectedArticle && (
        <div style={modalOverlay}>
          <div style={modalContent}>
            <h2>{selectedArticle.title}</h2>
            <p style={{ fontSize: "14px", color: "#4a5568", marginTop: "8px" }}>
              {selectedArticle.author}
            </p>
            <hr style={{ margin: "12px 0" }} />
            <p style={{ textAlign: "justify", lineHeight: "1.6" }}>
              {selectedArticle.content_long ||
                selectedArticle.content ||
                "Bu makale için detaylı içerik bulunamadı."}
            </p>
            <button onClick={() => setSelectedArticle(null)} style={closeBtn}>
              Kapat
            </button>
          </div>
        </div>
      )}
    </section>
  );
}

// Modal stilleri
const modalOverlay = {
  position: "fixed",
  top: 0,
  left: 0,
  width: "100%",
  height: "100%",
  backgroundColor: "rgba(0,0,0,0.6)",
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  zIndex: 1000,
};

const modalContent = {
  backgroundColor: "white",
  padding: "25px",
  borderRadius: "10px",
  width: "70%",
  maxHeight: "80vh",
  overflowY: "auto",
  boxShadow: "0 5px 20px rgba(0,0,0,0.2)",
};

const closeBtn = {
  marginTop: "20px",
  background: "#3b82f6",
  color: "white",
  border: "none",
  borderRadius: "6px",
  padding: "8px 18px",
  cursor: "pointer",
  fontWeight: "600",
};
