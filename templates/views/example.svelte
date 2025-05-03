<script>
    import {getContext} from "svelte";

    /**
     * Server data.
     * @typedef Data
     * @property {any} prop1
     */

    /** @type {Data} */
    const data = getContext("data")
</script>

<span>Hello, this is view!</span>
<!-- Content. -->