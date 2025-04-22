<style>
    a {
        width: 100%;
        cursor: default;
        border: 0;
        text-decoration: none;
        background: transparent;
    }

    a:hover {
        cursor: default;
        text-decoration: none;
    }

    .start {
        text-align: start;
    }

    .center {
        text-align: center;
    }

    .end {
        text-align: end;
    }
</style>

<script>
    import {getContext} from "svelte";

    /**
     * @typedef Props
     * @property {string} view
     * @property {import("svelte").Snippet} children
     * @property {"start"|"center"|"end"} [align]
     * @property {Record<string,string>} [parameters]
     */

    /** @type {Props} */
    const {
        view,
        children,
        align = "start",
        parameters = {},
        ...rest
    } = $props()
    const navigate = getContext("navigate")
    const path = getContext("path")

    /**
     * @param {Event} e
     */
    function onmouseup(e) {
        e.preventDefault()
        navigate(view, parameters)
    }
</script>

<a href="{path(view, parameters)}"
   class:start={"start"===align}
   class:center={"center"===align}
   class:end={"end"===align}
   {onmouseup}
   {...rest}
>
    {@render children()}
</a>