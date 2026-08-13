<script lang="ts">
    import { onMount } from "svelte";
    import {
        SvelteFlow,
        Controls,
        Background,
        BackgroundVariant,
        type Node,
        MiniMap,
    } from "@xyflow/svelte";
    import "@xyflow/svelte/dist/style.css";

    import MainArtworkNode from "../component/MainArtworkNode.svelte";
    import ReferenceNode from "../component/ReferenceNode.svelte";
    import type { BoardResponse } from "../types/KritaBoardTypes";
    import { API_URL } from "../common";
    import { ChevronLeft, RotateCw } from "@lucide/svelte";

    // --- Svelte 5 Runes ---
    let nodes = $state<Node[]>([]);
    let boardData = $state<BoardResponse | null>(null);
    let isLoading = $state(true);

    let workspaceId = $state<string>("");
    let boardId = $state<string>("");

    // Map custom node types
    const nodeTypes = {
        mainArtwork: MainArtworkNode,
        referenceImage: ReferenceNode,
    };

    async function boardDataInit() {
        try {
            workspaceId = window.location.pathname.split("/")[2];
            boardId = window.location.pathname.split("/")[4];
            const res = await fetch(
                `${API_URL}/krita/${workspaceId}/board/${boardId}`,
            );
            boardData = await res.json();

            if (boardData) {
                // Origin is top-left of main canvas centered at (-width/2, -height/2)
                const canvasOriginX = -boardData.canvas.width / 2;
                const canvasOriginY = -boardData.canvas.height / 2;
                // Convert Go backend items into Svelte Flow Nodes
                nodes = boardData.items.map((item) => {
                    const isMain = item.type === "main_artwork";

                    const posX = isMain
                        ? canvasOriginX
                        : canvasOriginX + item.x;
                    const posY = isMain
                        ? canvasOriginY
                        : canvasOriginY + item.y;

                    return {
                        id: item.id,
                        type: isMain ? "mainArtwork" : "referenceImage",
                        position: { x: posX, y: posY },
                        width: item.width,
                        height: item.height,
                        data: {
                            url: item.url,
                            width: item.width,
                            height: item.height,
                            opacity: item.opacity,
                            fileName: item.fileName,
                        },
                    };
                });
            }
        } catch (err) {
            console.error("Failed to load board data:", err);
        } finally {
            isLoading = false;
        }
    }

    async function refreshBoard() {
        isLoading = true;
        nodes = [];
        boardData = null;

        await boardDataInit();
    }

    onMount(async () => {
        await boardDataInit();
    });
</script>

<main class="viewport">
    {#if isLoading}
        <div class="loading">Parsing Krita structure...</div>
    {:else if boardData}
        <!-- Floating Info Badge -->

        <div class="hud-badge top-left">
            <div class="hud-top-bar">
                <a href={`/krita-workspace/${workspaceId}`} class="back-button">
                    <ChevronLeft size={24} />
                </a>
                <button class="refresh-button" onclick={refreshBoard}>
                    <RotateCw size={20} />
                </button>
            </div>
        </div>

        <div class="hud-badge top-right">
            <div class="hud-header">
                <span class="hud-title"
                    >{boardData.canvas.name || boardData.fileName}</span
                >
            </div>

            <div class="hud-meta">
                {boardData.canvas.width} × {boardData.canvas.height} px • {boardData
                    .items.length - 1} references
            </div>
        </div>

        <!-- Infinite Svelte Flow Board -->
        <SvelteFlow
            {nodes}
            {nodeTypes}
            fitView
            fitViewOptions={{ padding: 0.2 }}
            minZoom={0.01}
            maxZoom={10}
            colorMode="dark"
        >
            <MiniMap />
            <Background
                variant={BackgroundVariant.Dots}
                gap={280}
                size={20}
                bgColor="#27272a"
            />
            <Controls />
        </SvelteFlow>
    {/if}
</main>

<style>
    :global(body) {
        margin: 0;
        padding: 0;
        font-family:
            system-ui,
            -apple-system,
            sans-serif;
        overflow: hidden;
    }

    .viewport {
        width: 100vw;
        height: 100vh;
        position: relative;
    }

    .loading {
        display: flex;
        height: 100%;
        align-items: center;
        justify-content: center;
        color: #a1a1aa;
        font-size: 0.875rem;
    }

    /* HUD Overlay styling */
    .hud-badge {
        position: absolute;
        top: 1rem;
        left: 1rem;
        z-index: 10;
        background-color: rgba(24, 24, 27, 0.9);
        border: 1px solid #27272a;
        padding: 0.75rem 1rem;
        border-radius: 0.25rem;
        backdrop-filter: blur(12px);
        display: flex;
        flex-direction: column;
        gap: 0.375rem;
    }

    .hud-badge.top-left {
        left: 1rem;
        padding: 0.5em;
        flex-direction: row;
        align-items: center;
    }

    .hud-badge.top-right {
        left: auto;
        right: 1rem;
    }

    .hud-top-bar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.5rem;
        margin-bottom: 0.25rem;
    }

    /* Back Button Styling */
    .back-button,
    .refresh-button {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.25rem;
        color: #a1a1aa;
        text-decoration: none;
        font-size: 0.75rem;
        font-weight: 500;
        padding: 0.25rem;
        aspect-ratio: 1 / 1;
        width: 2rem;
        border: none;
        border-radius: 0.375rem;
        transition: all 0.15s ease;
        background: none;
        cursor: pointer;
    }

    .back-button:hover,
    .refresh-button:hover {
        color: #ffffff;
    }

    .hud-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .hud-title {
        font-weight: 600;
        font-size: 0.875rem;
        color: #f4f4f5;
    }

    .hud-meta {
        font-size: 0.75rem;
        color: #71717a;
    }

    /* Strip default Svelte Flow Node borders and backgrounds */
    :global(.svelte-flow__node) {
        padding: 0 !important;
        border: none !important;
        background: transparent !important;
        box-shadow: none !important;
    }
</style>
