<script lang="ts">
    //:app-imports
    type Props = {
        id: string
        data: Record<string, any>
        ids: Record<string, string>
    }
    let {id: idInput, data: dataInput, ids: idsInput}:Props = $props()
    let server = $state({id: idInput, data: dataInput, ids: idsInput})
</script>

<!--app-router-->