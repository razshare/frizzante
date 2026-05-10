import { svelte } from "@sveltejs/vite-plugin-svelte"
import path from "path"
import { fileURLToPath } from "url"
import { defineConfig } from "vite"

const file = fileURLToPath(import.meta.url)
const dir = path.dirname(file).replace(/\\+/, "/")
const prod = (process.env.PROD ?? "0") === "1"

let sourcemap: "inline" | boolean = false
if (!prod) {
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
            "$gen": `${path.resolve(dir, "../.gen")}`,
            "$exports.client": `${path.resolve(dir, "./exports.client.ts")}`,
            "$exports.server": `${path.resolve(dir, "./exports.server.ts")}`,
        },
    },
    build: {
        copyPublicDir: false,
        sourcemap,
        rollupOptions: {
            input: { index: "./index.html" },
            external: /\.\/dist\/\.*/,
        },
    },
})
