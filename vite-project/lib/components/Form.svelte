<script>
    import {getContext} from "svelte";
    import {update} from "../scripts/update.js";

    /** @type {function(string):(function(Record<string,any>):void)} */
    const findNavigateByPath = getContext("findNavigateByPath")
    const data = getContext("data")
    const onsubmit = update({findNavigateByPath, data})

    /** @type {function(string,Record<string,any>):string} */
    const findPathByPageName = getContext("findPathByPageName")
    /**
     * @typedef Props
     * @property {string} [page]
     * @property {Record<string,string>} [parameters]
     * @property {import("svelte").Snippet} children
     */

    /** @type {Props} */
    let {
        page = '',
        parameters = {},
        children,
        ...rest
    } = $props()
</script>

<form method="POST" action={findPathByPageName(page, parameters)} {...rest} {onsubmit}>
    {@render children()}
</form>
