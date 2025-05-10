<style>
    form {
        position: relative;
        width: 100%;
        height: 100%;
    }

    .submit {
        display: none;
    }
</style>

<script>
    import {update} from "../scripts/update.js";
    import {getContext} from "svelte";
    import {uuid} from "../scripts/uuid.js";

    /** @type {function(string):string} */
    const path = getContext("path")
    /** @type {function(string):{view:string,parameters:Record<string,string>}} */
    const loadView = getContext("view")
    /** @type {function(string,Record<string,string>,false|Record<string,any>)} */
    const navigate = getContext("navigate")
    /** @type {Record<string,any>} */
    const data = getContext("data")

    const onsubmit = update({view: loadView, navigate, data})
    const id = uuid()

    /**
     * @typedef Props
     * @property {import("svelte").Snippet} children
     * @property {string} [view]
     * @property {Record<string,string|number|boolean>} [form]
     */

    /** @type {Props} */
    let {
        view = '',
        children,
        form = {},
    } = $props()


    if ('' !== view) {
        view = path(view)
    }
</script>

<form method="POST" action={view} {onsubmit}>
    {#each Object.keys(form) as key}
        {@const value = form[key]}
        <input type="hidden" name="{key}" value="{value}">
    {/each}

    <input class="submit" type="submit" id="{id}"/>

    <label for="{id}">
        {@render children()}
    </label>
</form>
