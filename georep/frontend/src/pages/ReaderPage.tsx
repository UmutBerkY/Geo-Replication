import React, { useEffect, useState } from "react";
import { apiGet } from "../api";
import ArticleCard from "../components/ArticleCard";
import { Article, ReplicationStatus, Session } from "../types";

type Props = {
  session: Session;
  onLogout: () => void;
};

export default function ReaderPage({ session, onLogout }: Props) {
  const [articles, setArticles] = useState<Article[]>([]);
  const [status, setStatus] = useState<ReplicationStatus[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [loadingStatus, setLoadingStatus] = useState(false);

  // 🔁 LATENCY DURUMU
  const [latencyText, setLatencyText] = useState<string | null>(null);
  const [loadingLatency, setLoadingLatency] = useState(false);

  const loadArticles = async () => {
    setError(null);
    setLoading(true);
    try {
      const data = await apiGet<Article[]>(`/articles?region=${session.region}`);
      setArticles(data || []);
    } catch (e) {
      setError(
        e instanceof Error ? e.message : "Haberler yüklenirken hata oluştu."
      );
    } finally {
      setLoading(false);
    }
  };

  const loadStatus = async () => {
    setLoadingStatus(true);
    try {
      const s = await apiGet<ReplicationStatus[]>("/replication-status");

      // 🧠 Master olan EU replikasını filtrele
      const filtered = (s || []).filter(
        (rep) => rep.replica.toLowerCase() !== "eu"
      );
      setStatus(filtered);
    } catch {
      setStatus([]);
    } finally {
      setLoadingStatus(false);
    }
  };

  // 🔎 GECİKME ÖLÇÜMÜ – backend’in döndürdüğü stringi OLDUĞU GİBİ kullan
  const loadLatency = async () => {
    setLoadingLatency(true);
    try {
      const res = await apiGet<{
        region: string;
        latency: string; // ⏱ Master’a göre gecikme kazancı… şeklinde hazır metin
        measured: string;
      }>(`/latency?region=${session.region}`);

      // ❗ Burada sadece stringi kaydediyoruz, ekstra format YOK
      setLatencyText(res.latency);
    } catch (e) {
      setLatencyText("Gecikme ölçümü yapılamadı.");
    } finally {
      setLoadingLatency(false);
    }
  };

  // Sayfa açıldığında ve bölge değiştiğinde otomatik yükleme + periyodik polling
  useEffect(() => {
    loadArticles();
    loadStatus();
    loadLatency(); // Bölge seçildiğinde latency de ölçülsün

    const interval = setInterval(() => {
      loadStatus();
      loadArticles();
    }, 3000);

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

          {/* ⏱ LATENCY METNİ – SADECE BACKEND’İN GÖNDERDİĞİ METİN */}
          {latencyText && (
            <p
              className="hint"
              style={{ marginTop: "0.5rem", fontWeight: 500 }}
            >
              {latencyText}
            </p>
          )}
        </div>
        <button onClick={onLogout}>Çıkış</button>
      </div>

      {error && <p className="error">{error}</p>}

      <div className="actions" style={{ gap: "0.5rem" }}>
        <button onClick={loadArticles} disabled={loading}>
          {loading ? "Yükleniyor..." : "Haberleri yenile"}
        </button>
        <button onClick={loadStatus} disabled={loadingStatus}>
          {loadingStatus ? "Durum getiriliyor..." : "Replikasyon durumu"}
        </button>
        <button onClick={loadLatency} disabled={loadingLatency}>
          {loadingLatency ? "Gecikme ölçülüyor..." : "Gecikme kazancını ölç"}
        </button>
      </div>

      <div className="stories">
        {articles.length > 0 ? (
          articles.map((a) => <ArticleCard key={a.id} article={a} />)
        ) : (
          <p className="hint">
            Bu replikada henüz veri yok. Bir yazar EU master üzerinden yeni
            haber yayınladığında burada görünecek.
          </p>
        )}
      </div>

      {/* Replikasyon Durumu (EU hariç) */}
      {status.length > 0 && (
        <div
          style={{
            backgroundColor: "#f0f9ff",
            border: "1px solid #bae6fd",
            borderRadius: "8px",
            padding: "15px",
            marginTop: "1rem",
          }}
        >
          <h3 style={{ marginTop: 0, color: "#1e40af" }}>
            🔄 Replikasyon Durumu
          </h3>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))",
              gap: "10px",
            }}
          >
            {status.map((s, i) => {
              const colorMap =
                {
                  ok: {
                    bg: "#c6f6d5",
                    text: "#22543d",
                    border: "#9ae6b4",
                    icon: "✅",
                  },
                  error: {
                    bg: "#fed7d7",
                    text: "#742a2a",
                    border: "#fc8181",
                    icon: "❌",
                  },
                  syncing: {
                    bg: "#feebc8",
                    text: "#744210",
                    border: "#fbd38d",
                    icon: "🔄",
                  },
                }[s.status] || {
                  bg: "#edf2f7",
                  text: "#2d3748",
                  border: "#e2e8f0",
                  icon: "ℹ️",
                };

              return (
                <div
                  key={i}
                  style={{
                    backgroundColor: colorMap.bg,
                    color: colorMap.text,
                    padding: "10px",
                    borderRadius: "6px",
                    border: `1px solid ${colorMap.border}`,
                    fontSize: "14px",
                    fontWeight: "500",
                  }}
                >
                  {colorMap.icon} {s.replica}: {s.status}
                  {s.last_at && (
                    <div
                      style={{
                        fontSize: "12px",
                        marginTop: "4px",
                        opacity: 0.8,
                      }}
                    >
                      {new Date(s.last_at).toLocaleTimeString("tr-TR")}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      )}
    </section>
  );
}
