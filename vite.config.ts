import vue from "@vitejs/plugin-vue";

import path from "path";
import { fileURLToPath } from "node:url";

import { defineConfig, loadEnv } from "vite";
import vueDevTools from "vite-plugin-vue-devtools";
import { quasar, transformAssetUrls } from "@quasar/vite-plugin";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const backendUrl = env.VITE_BACKEND_URL || "http://localhost:8090";
  return {
    plugins: [
      vue({ template: { transformAssetUrls } }),
      // @quasar/plugin-vite options list:
      // https://github.com/quasarframework/quasar/blob/dev/vite-plugin/index.d.ts
      quasar({
        sassVariables: fileURLToPath(
          new URL("./frontend/src/quasar-variables.sass", import.meta.url)
        ),
      }),
      vueDevTools(),
    ],
    build: {
      outDir: "./backend/dist",
      emptyOutDir: true,
      chunkSizeWarningLimit: 1600,
    },
    publicDir: "./frontend/public",
    resolve: { alias: { "@": path.resolve(__dirname, "./frontend/src") } },
    server: {
      proxy: {
        "/api": {
          target: backendUrl,
          changeOrigin: true,
        },
        "/maps": {
          target: backendUrl,
          changeOrigin: true,
        },
      },
    },
  };
});
