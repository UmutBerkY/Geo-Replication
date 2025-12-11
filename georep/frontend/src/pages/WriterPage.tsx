import React, { useEffect, useState } from "react";
import { apiGet, apiPost } from "../api";
import { Article, ReplicationStatus, Session } from "../types";

type Props = {
  session: Session;
  onLogout: () => void;
  onArticleAdded?: () => void;
};

export default function WriterPage({ session, onLogout, onArticleAdded }: Props) {
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState(session.username || "");
  const [summary, setSummary] = useState("");
  const [contentLong, setContentLong] = useState("");
  const [loading, setLoading] = useState(false);
  const [articles, setArticles] = useState<Article[]>([]);
  const [repStatus, setRepStatus] = useState<ReplicationStatus[]>([]);
  const [expandedId, setExpandedId] = useState<number | null>(null);

  const loadArticles = async () => {
    const data = await apiGet<Article[]>("/articles?region=eu");
    setArticles(data || []);
  };

  const loadRepStatus = async () => {
    const s = await apiGet<ReplicationStatus[]>("/replication-status");
    setRepStatus((s || []).filter((r) => r.replica.toLowerCase() !== "eu"));
  };

  useEffect(() => {
    loadArticles();
    loadRepStatus();
    const interval = setInterval(loadRepStatus, 3000);
    return () => clearInterval(interval);
  }, []);

  const publish = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const created = await apiPost<Article>("/articles", {
        title,
        summary,
        content: summary,
        content_long: contentLong,
        author,
      });
      setArticles((prev) => [created, ...prev]);
      setTitle("");
      setSummary("");
      setContentLong("");
      setAuthor(session.username || "");
      if (onArticleAdded) onArticleAdded();
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Bu haberi silmek istediğine emin misin?")) return;
    try {
      await fetch(
        `${
          import.meta.env.VITE_API_BASE || "http://localhost:8080/api"
        }/articles/${id}`,
        { method: "DELETE" }
      );
      await loadArticles();
    } catch {
      alert("Silme işlemi başarısız oldu.");
    }
  };

  const toggleExpand = (id: number) => {
    setExpandedId((prev) => (prev === id ? null : id));
  };

  return (
    <section className="card">
      <div className="header-row">
        <div>
          <h2>Yazar Görünümü – EU Master</h2>
          <p className="hint">
            Tüm yazılar EU master’a kaydedilir. Diğer replikalar senkronize olur.
          </p>
        </div>
        <button onClick={onLogout}>Çıkış</button>
      </div>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "1fr 1fr",
          gap: "20px",
          marginTop: "1.5rem",
        }}
      >
        {/* 🔄 Replikasyon Durumu */}
        <div
          style={{
            backgroundColor: "#f0f9ff",
            border: "1px solid #bae6fd",
            borderRadius: "8px",
            padding: "20px",
          }}
        >
          <h3 style={{ marginTop: 0, color: "#1e40af" }}>🔄 Replikasyon Durumu</h3>
          {repStatus.map((s, i) => {
            const colorMap =
              {
                ok: {
                  bg: "#c6f6d5",
                  text: "#22543d",
                  border: "#9ae6b4",
                  icon: "✅",
                },
                syncing: {
                  bg: "#feebc8",
                  text: "#744210",
                  border: "#fbd38d",
                  icon: "🔄",
                },
                error: {
                  bg: "#fed7d7",
                  text: "#742a2a",
                  border: "#fc8181",
                  icon: "❌",
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
              </div>
            );
          })}
        </div>

        {/* 📝 Yeni Makale Ekle */}
        <div
          style={{
            backgroundColor: "white",
            border: "1px solid #e2e8f0",
            borderRadius: "8px",
            padding: "20px",
          }}
        >
          <h3 style={{ marginTop: 0, color: "#2d3748" }}>📝 Yeni Makale Ekle</h3>
          <form
            onSubmit={publish}
            style={{ display: "flex", flexDirection: "column", gap: "12px" }}
          >
            <input
              placeholder="Haber başlığı"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              style={{
                padding: "10px",
                borderRadius: "6px",
                border: "1px solid #cbd5e0",
              }}
            />
            <input
              placeholder="Yazar adı"
              value={author}
              onChange={(e) => setAuthor(e.target.value)}
              required
              style={{
                padding: "10px",
                borderRadius: "6px",
                border: "1px solid #cbd5e0",
              }}
            />
            <textarea
              placeholder="Kısa özet (önizleme için)"
              value={summary}
              onChange={(e) => setSummary(e.target.value)}
              required
              rows={2}
              style={{
                padding: "10px",
                borderRadius: "6px",
                border: "1px solid #cbd5e0",
                resize: "vertical",
              }}
            />
            <textarea
              placeholder="Uzun içerik (makalenin tamamı)"
              value={contentLong}
              onChange={(e) => setContentLong(e.target.value)}
              required
              rows={8}
              style={{
                padding: "10px",
                borderRadius: "6px",
                border: "1px solid #cbd5e0",
                resize: "vertical",
              }}
            />
            <button
              type="submit"
              disabled={loading}
              style={{
                backgroundColor: loading ? "#9ca3af" : "#3b82f6",
                color: "white",
                border: "none",
                padding: "10px",
                borderRadius: "6px",
                fontWeight: "600",
                cursor: "pointer",
              }}
            >
              {loading ? "Yayınlanıyor..." : "Yayınla"}
            </button>
          </form>
        </div>
      </div>

      <hr style={{ margin: "2rem 0 1rem 0" }} />
      <h3>EU Master’daki Yazılar</h3>
      <div className="stories">
        {articles.length > 0 ? (
          articles.map((a) => (
            <div
              key={a.id}
              style={{
                backgroundColor: "#fff",
                border: "1px solid #e2e8f0",
                borderRadius: "8px",
                padding: "16px",
                marginBottom: "16px",
              }}
            >
              <h4 style={{ margin: "0 0 4px 0" }}>{a.title}</h4>
              <p style={{ margin: "0 0 8px 0", color: "#4b5563" }}>
                <strong>{a.author}</strong> — {new Date(a.created_at).toLocaleString()}
              </p>
              <p style={{ color: "#374151" }}>
                {expandedId === a.id ? a.content_long : a.summary}
              </p>
              <div style={{ marginTop: "10px", display: "flex", gap: "10px" }}>
                <button
                  onClick={() => toggleExpand(a.id)}
                  style={{
                    backgroundColor: "#3b82f6",
                    color: "white",
                    border: "none",
                    padding: "6px 12px",
                    borderRadius: "6px",
                    fontSize: "14px",
                    cursor: "pointer",
                  }}
                >
                  {expandedId === a.id ? "Kısalt" : "Devamını Oku"}
                </button>
                <button
                  onClick={() => handleDelete(a.id)}
                  style={{
                    backgroundColor: "#ef4444",
                    color: "white",
                    border: "none",
                    padding: "6px 12px",
                    borderRadius: "6px",
                    fontSize: "14px",
                    cursor: "pointer",
                  }}
                >
                  Sil
                </button>
              </div>
            </div>
          ))
        ) : (
          <p className="hint">Henüz makale yok.</p>
        )}
      </div>
    </section>
  );
}
