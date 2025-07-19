import { defineConfig } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"

const IS_DEV = "1" === (process.env.DEV ?? "")

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
            $lib: "./lib",
            $frizzante: "./frizzante",
            $exports: "./exports",
        },
    },
    build: {
        sourcemap: IS_DEV ? "inline" : false,
        rollupOptions: {
            input: {
                index: "./index.html",
            },
        },
    },
})
