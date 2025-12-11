import { useState, FormEvent } from "react";
import { apiPost } from "../api";
import { Region, Role, Session } from "../types";

type Props = {
  onLogin: (s: Session) => void;
};

const regions: Region[] = ["eu", "us", "asia", "sa", "tr", "africa"];

export default function LoginPage({ onLogin }: Props) {
  const [username, setUsername] = useState("");
  const [role, setRole] = useState<Role>("reader");
  const [region, setRegion] = useState<Region>("eu");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const submit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const resp = await apiPost<Session>("/login", {
        username,
        role,
        region
      });
      onLogin(resp);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Giriş sırasında bir hata oluştu."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="card">
      <h2>Giriş yap ve rol/bölge seç</h2>
      <form onSubmit={submit} className="form">
        <input
          placeholder="Kullanıcı adı"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <label>
          Rol:
          <select value={role} onChange={(e) => setRole(e.target.value as Role)}>
            <option value="reader">Okuyucu</option>
            <option value="writer">Yazıcı</option>
          </select>
        </label>
        <label>
          Bölge:
          <select
            value={region}
            onChange={(e) => setRegion(e.target.value as Region)}
          >
            {regions.map((r) => (
              <option key={r} value={r}>
                {r.toUpperCase()}
              </option>
            ))}
          </select>
        </label>
        <button type="submit" disabled={loading}>
          {loading ? "Giriş yapılıyor..." : "Devam et"}
        </button>
      </form>
      {error && <p className="error">{error}</p>}
      <p className="hint">
        Yazıcı her zaman EU master’a yazar. Okuyucu, seçtiği bölgedeki en yakın
        replikadan okur.
      </p>
    </section>
  );
}


