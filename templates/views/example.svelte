<script lang="ts">
    import {getContext} from "svelte";
    import type {ServerContext} from "$frizzante/types.ts";

    const server = getContext("server") as ServerContext<any>
</script>

<span>Hello, this is view!</span>
<!-- Content. -->