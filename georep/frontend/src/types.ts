export type Region = "eu" | "us" | "asia" | "sa" | "tr" | "africa";

export type Role = "reader" | "writer";

export type Session = {
  username: string;
  role: Role;
  region: Region;
  token: string;
};

export type Article = {
  id: number;
  title: string;
  content: string;
  author: string;
  region: string;
  created_at: string;
};

export type ReplicationStatus = {
  replica: string;
  status: string;
  last_at?: string;
};


