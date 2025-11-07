<style>
    svg {
        --size: 1rem;
        width: var(--size);
        height: var(--size);
    }
</style>

<script lang="ts">
    type Props = {
        path: string
        size?: string
        fill?: string
        stroke?: string
    }

    let { path, size = "1rem", fill = "currentColor", stroke = "none" }: Props = $props()
</script>

<svg xmlns="http://www.w3.org/2000/svg" style="--size:{size}" viewBox="0 0 24 24" {stroke} {fill}
    ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d={path} /></svg
>
