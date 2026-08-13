<script lang="ts">
    import { onMount } from "svelte";
    import { API_HOST, API_URL } from "../common";
    import KritaArtworkItemComponent from "../component/KritaArtworkItemComponent.svelte";
    import Breadcrumb from "../component/Breadcrumb.svelte";

    type Artwork = {
        name: string;
        filename: string;
        routePath: string;
        previewPath: string;
        createdAt: string; // format: 2026-08-10T20:01:33+07:00"
    };

    type WorkspaceItem = {
        name: string;
        routePath: string;
        fileCount: number;
        artworks: Artwork[];
    };

    let workspaceItems = $state<WorkspaceItem | null>(null);
    let loading = $state<boolean>(true);
    let error = $state<string | null>(null);
    let workspaceByPath = $state<string | undefined | null>(null);

    async function fetchWorkspaceItems(workspace: string) {
        try {
            const response = await fetch(`${API_URL}/krita/${workspace}`);
            const data = await response.json();
            workspaceItems = data as WorkspaceItem;
        } catch (err) {
            error = (err as Error).message;
            console.error(err);
        } finally {
            loading = false;
        }
    }

    onMount(async () => {
        workspaceByPath = window.location.pathname
            .split("/")
            .filter(Boolean)
            .pop();
        if (!workspaceByPath) return;
        await fetchWorkspaceItems(workspaceByPath);
    });

    const breadcrumbItems = $derived([
        {
            label: "Home",
            href: "/",
        },
        {
            label: "Krita Workspace",
            href: "/krita-workspace",
        },
        {
            label: workspaceByPath ?? "",
            href: window.location.pathname,
        },
    ]);
</script>

<section class="krita-workspace-item-panel">
    <Breadcrumb items={breadcrumbItems} />

    <section class="krita-workspace-container">
        {#if loading}
            <p>Loading...</p>
        {:else if error}
            <p>{error}</p>
        {:else if workspaceItems !== null}
            <ul class="krita-workspace-list">
                {#each workspaceItems.artworks as artwork}
                    <KritaArtworkItemComponent
                        preview={API_HOST + artwork.previewPath}
                        projectName={workspaceItems.name}
                        artworkName={artwork.name}
                        createdAt={artwork.createdAt}
                        description="Krita workspace project"
                    />
                {/each}
            </ul>
        {/if}
    </section>
</section>

<style>
    .krita-workspace-item-panel {
        position: relative;
        width: 100%;
        height: 100%;
        padding: 1.75em;
        box-sizing: border-box;
        padding-bottom: 0;

        display: flex;
        flex-direction: column;
        gap: 1.7em;
    }

    .krita-workspace-container {
        position: relative;
        width: 100%;
        height: 100%;

        overflow-y: auto;
        padding-bottom: 10em;
        box-sizing: border-box;
    }

    .krita-workspace-list {
        list-style: none;
        padding: 0;
        margin: 0;

        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(12em, 1fr));
        gap: 1.6em;
    }
</style>
