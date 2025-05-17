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
</style>

<script>
    import {getContext} from "svelte";

    /** @type {function(string,Record<string,any>):string} */
    const findPathByPageName = getContext("findPathByPageName")

    /** @type {function(string,Record<string,any>):void} */
    const navigate = getContext("navigate")

    /**
     * @typedef Props
     * @property {string} page
     * @property {"start"|"center"|"end"} [align]
     * @property {Record<string,string>} [parameters]
     * @property {import("svelte").Snippet} children
     */

    /** @type {Props} */
    const {
        page,
        parameters = {},
        children,
        ...rest
    } = $props()

    /**
     * @param {Event} e
     */
    function onmouseup(e) {
        e.preventDefault()
        navigate(page, parameters)
    }
</script>

<a href="{findPathByPageName(page, parameters)}" {onmouseup} {...rest}>
    {@render children()}
</a>