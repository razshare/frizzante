<script lang="ts">
    //:app-imports
    import {setContext} from "svelte"

    type Props = {
        view: string
        data: any
    }

    let {view: viewInput, data: dataInput, error: errorInput}: any = $props()
    const server = $state({view: viewInput, data: dataInput, error: errorInput}) as Props
    setContext("server", server)
</script>

<!--app-router-->