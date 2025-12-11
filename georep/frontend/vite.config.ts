import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

// Vite config for React frontend.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    host: "0.0.0.0"
  }
});



