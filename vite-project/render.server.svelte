<script lang="ts">
    //:app-imports
    import {setContext} from "svelte"

    type Props = {
        id: string
        data: any
        ids: Record<string, string>
    }

    let {id: idInput, data: dataInput, ids: idsInput}: any = $props()
    let server = $state({id: idInput, data: dataInput, ids: idsInput}) as Props
    setContext("server", server)
</script>

<!--app-router-->