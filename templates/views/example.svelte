<script>
    /**
     * @typedef ServerData
     * @property {any} prop1
     */
    /** @type {ServerProperties<ServerData>} */
    let {server = $bindable()} = $props()
</script>

<span>Hello, this is view!</span>
<!-- Content. -->