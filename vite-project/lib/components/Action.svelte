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
    import {action} from "$frizzante/scripts/action.ts";
    import {uuid} from "../scripts/uuid.ts";
    const id  = uuid()
    type Props = {
        of: string
        using?: any
        children: Snippet
    }

    let {
        of,
        using,
        children,
    }: Props = $props()
</script>

<form {...action(of)}>
    {#each Object.keys(using ?? {}) as key}
        {@const value = using[key]}
        <input type="hidden" name="{key}" value="{value}">
    {/each}

    <input class="submit" type="submit" id="{id}"/>

    <label for="{id}">
        {@render children()}
    </label>
</form>
