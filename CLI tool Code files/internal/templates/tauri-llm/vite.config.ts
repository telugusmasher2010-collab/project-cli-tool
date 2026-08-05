import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],

  // Vite options tailored for Tauri development.
  clearScreen: false,
  server: {
    port: 1420,
    strictPort: true,
    watch: {
      // Tauri is a Rust build; do not watch it.
      ignored: ["**/src-tauri/**"],
    },
  },
});
