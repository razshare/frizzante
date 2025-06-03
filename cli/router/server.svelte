<script lang="ts">
    //:imports
    import {setContext} from "svelte"

    let {name, data, error, renderMode} = $props()
    const view = $state({name, data, error, renderMode})
    setContext("view", view)
</script>

<!--:router-->