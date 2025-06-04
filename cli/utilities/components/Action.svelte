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

<script lang="ts">
    import type {Snippet} from "svelte";
    import {action} from "../scripts/action.ts";
    import {uuid} from "../scripts/uuid.ts";

    const id = uuid()
    type Props = {
        path: string
        using?: Record<string, unknown>
        children: Snippet
    }

    let {
        path,
        using,
        children,
    }: Props = $props()
</script>

<form {...action(path)}>
    {#each Object.keys(using ?? {}) as key(key)}
        {@const value = (using??{})[key]??''}
        <input type="hidden" name="{key}" value="{value}">
    {/each}

    <input class="submit" type="submit" id="{id}"/>

    <label for="{id}">
        {@render children()}
    </label>
</form>
