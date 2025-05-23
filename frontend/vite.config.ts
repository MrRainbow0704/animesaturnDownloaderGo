import {defineConfig} from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      "$src": "/src",
      "$assets": "/src/assets",
      "$lib": "/src/lib",
      "$components": "/src/components",
      "$wails": "/wailsjs"
    }
  }
})
