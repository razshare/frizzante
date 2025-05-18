<script>
    //:app-imports

    /**
     * @typedef Props
     * @property {string} id
     * @property {Record<string,any>} data
     * @property {Record<string,string>} ids
     */

    /** @type {Props} */
    let {id: idInput, data: dataInput, ids: idsInput} = $props()
    let server = $state({id: idInput, data: dataInput, ids: idsInput})
</script>

<!--app-router-->