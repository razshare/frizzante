<script lang="ts">
    type Props = {
        path: string
        size?: string
        fill?: string
        stroke?: string
    }

    let { path, size = "1.6rem", fill = "currentColor", stroke = "none" }: Props = $props()
</script>

<svg xmlns="http://www.w3.org/2000/svg" style="width:{size};height:{size}" viewBox="0 0 24 24" {stroke} {fill}
    ><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d={path} /></svg
>
