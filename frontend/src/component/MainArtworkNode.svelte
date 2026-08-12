<script lang="ts">
    import type { NodeProps } from "@xyflow/svelte";
    import { API_HOST } from "../common";

    let { data }: NodeProps = $props();
</script>

<div
    class="main-art-node"
    style="width: {data.width}px; height: {data.height}px;"
>
    <img src={API_HOST + data.url} alt="Main Krita Artwork" />
</div>

<style>
    .main-art-node {
        border: 2px solid #3b82f6;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.6);
        background-color: #18181b;
    }

    .main-art-node img {
        width: 100%;
        height: 100%;
        object-fit: contain;
        display: block;
    }
</style>
