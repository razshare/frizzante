import {defineConfig} from "vite";
import {svelte} from "@sveltejs/vite-plugin-svelte";

let sourcemap: false | "inline" = false;

if ("1" === (process.env.DEV ?? "")) {
    sourcemap = "inline";
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
            "$client": "./app/client.ts",
            "$server": "./app/server.ts",
            "$frz": "./app/frz",
            "$lib": "./app/lib",
        },
    },
    build: {
        sourcemap,
        rollupOptions: {
            input: {
                index: "./index.html",
            },
        },
    },
});
