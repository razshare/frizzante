<script>
    import {getContext} from "svelte";
    import {update} from "../scripts/update.js";

    /** @type {function(string):string} */
    const path = getContext("path")
    /** @type {function(string):{view:string,parameters:Record<string,string>}} */
    const loadView = getContext("view")
    /** @type {function(string,Record<string,string>,false|Record<string,any>)} */
    const navigate = getContext("navigate")
    /** @type {Record<string,any>} */
    const data = getContext("data")

    const onsubmit = update({view: loadView, navigate, data})

    /**
     * @typedef Props
     * @property {import("svelte").Snippet} children
     * @property {string} [view]
     */

    /** @type {Props} */
    let {children, view = '', ...rest} = $props()

    if ('' !== view) {
        view = path(view)
    }

</script>

<form method="POST" action={view} {...rest} {onsubmit}>
    {@render children()}
</form>
