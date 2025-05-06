<script>
    import {getContext} from "svelte";
    import {update} from "../scripts/update.js";

    /** @type {function(string):string} */
    const path = getContext("path")
    /** @type {function(string):{view:string,parameters:Record<string,string>}} */
    const view = getContext("view")
    /** @type {function(string,Record<string,string>,false|Record<string,any>):void} */
    const navigate = getContext("navigate")
    /** @type {Record<string,any>} */
    const data = getContext("data")

    const onsubmit = update({view, navigate, data})

    /**
     * @typedef Props
     * @property {import("svelte").Snippet} children
     * @property {string} [action]
     */

    /** @type {Props} */
    let {children, action = '?', ...rest} = $props()

    if ('?' !== action) {
        action = path(action)
    }

</script>

<form method="POST" {action} {...rest} {onsubmit}>
    {@render children()}
</form>
