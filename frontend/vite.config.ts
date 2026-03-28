import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // Allow requests proxied from the Go backend in dev mode
    allowedHosts: "all",
    proxy: {
      "/api/terminal": {
        target: "http://localhost:8000",
        ws: true,
        changeOrigin: true,
        timeout: 0,
        proxyTimeout: 0,
      },
      "/api": {
        target: "http://localhost:8000",
        changeOrigin: true,
      },
      "/ws": {
        target: "http://localhost:8000",
        ws: true,
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
  },
});
