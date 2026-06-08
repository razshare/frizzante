import { svelte } from "@sveltejs/vite-plugin-svelte"
import path from "path"
import { fileURLToPath } from "url"
import { defineConfig } from "vite"

const file = fileURLToPath(import.meta.url)
const dir = path.dirname(file).replace(/\\+/, "/")
const dev = (process.env.DEV ?? "0") === "1"

let sourcemap: "inline" | boolean = false
if (dev) {
    sourcemap = "inline"
}

// https://vite.dev/config/
export default defineConfig({
    plugins: [
        svelte({
            compilerOptions: {
                css: "injected",
            },
        }),
    ],
    resolve: {
        alias: {
            "$lib": `${path.resolve(dir, "./lib")}`,
            "$exports.client": `${path.resolve(dir, "./exports.client.ts")}`,
            "$exports.server": `${path.resolve(dir, "./exports.server.ts")}`,
        },
    },
    build: {
        sourcemap,
        rollupOptions: {
            input: { index: "./index.html" },
            external: /\.\/dist\/\.*/,
        },
    },
})
