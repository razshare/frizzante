<script lang="ts">
    import {getContext, onMount} from "svelte";
    import {route} from "$frizzante/scripts/route.ts";
    import type {ServerContext} from "$frizzante/types.ts";

    const server = getContext("server") as ServerContext<any>
    onMount(function () {
        route(server)
    })
</script>