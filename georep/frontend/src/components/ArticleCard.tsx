import { Article } from "../types";

type Props = {
  article: Article;
  onDelete?: (id: number) => void; // 🔹 Silme için opsiyonel callback
};

export default function ArticleCard({ article, onDelete }: Props) {
  return (
    <article className="story">
      <h3>{article.title}</h3>
      <p className="story-meta">
        Yazar: {article.author} • Bölge: {article.region.toUpperCase()} •{" "}
        {new Date(article.created_at).toLocaleString()}
      </p>
      <p className="story-body">{article.content}</p>

      {onDelete && (
        <button
          onClick={() => onDelete(article.id)}
          style={{
            marginTop: "0.5rem",
            backgroundColor: "#c0392b",
            color: "white",
            border: "none",
            padding: "6px 12px",
            borderRadius: "5px",
            cursor: "pointer",
          }}
        >
          Sil
        </button>
      )}
    </article>
  );
}
