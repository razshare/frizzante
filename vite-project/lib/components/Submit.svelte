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
    import {getContext} from "svelte";
    import {update} from "../scripts/update.js";
    import {uuid} from "../scripts/uuid.js";

    /** @type {function(string):(function(Record<string,any>):void)} */
    const findNavigateByPath = getContext("findNavigateByPath")
    const data = getContext("data")
    const onsubmit = update({findNavigateByPath, data})

    /** @type {function(string,Record<string,any>):string} */
    const findPathByPageName = getContext("findPathByPageName")
    const id = uuid()

    /**
     * @typedef Props
     * @property {string} [page]
     * @property {Record<string,string|number|boolean>} [form]
     * @property {Record<string,string>} [parameters]
     * @property {import("svelte").Snippet} children
     */

    /** @type {Props} */
    let {
        page = '',
        parameters = {},
        form = {},
        children,
    } = $props()
</script>

<form method="POST" action={findPathByPageName(page, parameters)} {onsubmit}>
    {#each Object.keys(form) as key}
        {@const value = form[key]}
        <input type="hidden" name="{key}" value="{value}">
    {/each}

    <input class="submit" type="submit" id="{id}"/>

    <label for="{id}">
        {@render children()}
    </label>
</form>
