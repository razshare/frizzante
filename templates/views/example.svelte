<script lang="ts">
    import {getContext} from "svelte";
    import type {ServerContext} from "$frizzante/types.ts";
    import Router from "$frizzante/components/Router.svelte";

    const server = getContext("server") as ServerContext<any>
</script>

<Router/>
<span>Hello, this is view!</span>
<!-- Content. -->