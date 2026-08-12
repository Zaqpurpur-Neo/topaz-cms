<script lang="ts">
    import type { NodeProps } from "@xyflow/svelte";
    import { API_HOST } from "../common";

    let { data }: NodeProps = $props();
</script>

<div
    class="reference-node"
    style="
    width: {data.width}px;
    height: {data.height}px;
    opacity: {data.opacity};
  "
>
    <img src={API_HOST + data.url} alt={data.fileName} />
    <span class="ref-label">{data.fileName}</span>
</div>

<style>
    .reference-node {
        position: relative;
        border: 1px solid #3f3f46;
        border-radius: 4px;
        overflow: hidden;
        background-color: #18181b;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.4);
        transition: border-color 0.15s ease;
    }

    .reference-node:hover {
        border-color: #60a5fa;
    }

    .reference-node img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        display: block;
    }

    .ref-label {
        position: absolute;
        bottom: 4px;
        left: 4px;
        opacity: 0;
        background-color: rgba(0, 0, 0, 0.75);
        color: #d4d4d8;
        font-size: 10px;
        padding: 2px 4px;
        border-radius: 2px;
        transition: opacity 0.15s ease;
    }

    .reference-node:hover .ref-label {
        opacity: 1;
    }
</style>
