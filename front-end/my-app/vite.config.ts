// vite.config.ts
import { defineConfig } from 'vite';
import angular from '@analogjs/vite-plugin-angular';

export default defineConfig({
  plugins: [
    angular(),       // active la compilation Angular sous Vite
  ],
  server: {
    port: 4200,
    proxy: {
      // proxy toutes les requêtes /api/* vers ton back Go sur 8080
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        // retire le préfixe /api avant de forwarder
        rewrite: path => path.replace(/^\/api/, '')
      }
    }
  }
});
