<script lang="ts">
    import { onMount } from "svelte";
    import { API_URL } from "../common";
    import KritaWorkspaceComponent from "../component/KritaWorkspaceComponent.svelte";
    import Breadcrumb from "../component/Breadcrumb.svelte";

    const apiUrl = API_URL;

    type WorkspacesResponse = {
        workspaces: {
            name: string;
            routePath: string;
            fileCount: number;
            [key: string]: any;
        }[];
        workspacePath: string;
    };

    let workspacesResponse = $state<WorkspacesResponse | null>(null);
    let loading = $state<boolean>(true);
    let error = $state<string | null>(null);

    async function fetchWorkspaces() {
        try {
            const response = await fetch(`${apiUrl}/krita`);
            const data = await response.json();
            return data;
        } catch (err) {
            error = err instanceof Error ? err.message : String(err);
            console.error(err);
        } finally {
            loading = false;
        }
    }

    onMount(async () => {
        workspacesResponse = await fetchWorkspaces();
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
    ]);
</script>

<section class="krita-workspace-panel">
    <Breadcrumb items={breadcrumbItems} />

    <section class="krita-workspace-container">
        {#if loading}
            <p>Loading...</p>
        {:else if error}
            <p>{error}</p>
        {:else if workspacesResponse !== null}
            {#if workspacesResponse.workspaces}
                <ul class="krita-workspace-list">
                    {#each workspacesResponse.workspaces as workspaceItem}
                        <KritaWorkspaceComponent
                            projectName={workspaceItem.name}
                            description="Krita workspace project"
                        />
                    {/each}
                </ul>
            {/if}
        {/if}
    </section>
</section>

<style>
    .krita-workspace-panel {
        position: relative;
        width: 100%;
        height: 100%;
        padding: 1.75em;
        box-sizing: border-box;

        display: flex;
        flex-direction: column;
        gap: 1.7em;
        padding-bottom: 0;
    }

    .krita-workspace-container {
        position: relative;
        width: 100%;
        height: 100%;

        overflow-y: auto;
        padding-bottom: 20em;
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
    .title {
        font-size: 1.5rem;
        font-weight: 600;
        color: #fff;
    }
</style>
