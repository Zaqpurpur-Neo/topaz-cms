<script lang="ts">
    import { Gem } from "@lucide/svelte";

    import {
        routerName,
        p,
        route,
        isActive,
        isSidebarClosed,
        getActiveRouteItem,
        navigate,
    } from "../router";
    import { isActiveLink } from "sv-router";

    import type { RouteItem } from "../router";
    import { get } from "svelte/store";

    let closedSidebar = $derived(isSidebarClosed());
    let currentRoute = $derived(getActiveRouteItem());

    function handleRouteClick(e: MouseEvent, item: RouteItem) {
        if (item.keepQueryParam && item.getQuery) {
            let query = item.getQuery();

            if (query) {
                e.preventDefault();

                query = query.replace(/^%3F/i, "?");
                if (!query.startsWith("?")) {
                    query = `?${query}`;
                }

                const basePath = p(item.path);
                const targetUrl = basePath + query;

                window.history.pushState(null, "", targetUrl);
                window.dispatchEvent(new PopStateEvent("popstate"));
            }
        }
    }
</script>

{#if !closedSidebar}
    <section class="side-panel">
        <div class="side-panel-header">
            <div class="icon-section">
                <Gem size={18} />
            </div>
            <div class="text-section">
                <h5>Topaz Workspace</h5>
                <p>Does all my work</p>
            </div>
        </div>

        <div class="gap-section"></div>

        <p class="side-label-title">Tools</p>
        <div class="route-side-panel">
            {#each routerName as route}
                {#if route.isVisible}
                    <a
                        class="route-link"
                        {@attach isActiveLink({
                            className: "route-link-active",
                            startsWith: false,
                        })}
                        onclick={(e) => handleRouteClick(e, route)}
                        href={p(route.path)}
                    >
                        <svelte:component this={route.icon} size={18} />
                        <span>{route.title}</span>
                    </a>
                {/if}
            {/each}
        </div>

        <div class="badge-side-panel">
            <div class="badge-wrapper">
                <img src="/topaz.webp" alt="topaz badge" />
            </div>
        </div>
    </section>
{/if}

<style>
    :global(.side-panel.close-sidebar) {
        display: none;
    }

    .side-panel {
        background-color: var(--gray-color-1);
        padding: 0.75em;
        min-width: 14em;
        width: min-content;

        display: flex;
        flex-direction: column;
        height: 100%;
        border-right: 1px solid #262626;
    }

    .gap-section {
        position: relative;
        height: 2em;
    }

    .side-panel-header {
        display: flex;
        align-items: center;

        gap: 0.5em;
        padding: 0 0.35em;
    }

    .icon-section {
        color: var(--text-color);
        background-color: var(--accent-red);
        padding: 0.4em;

        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 8px;
    }

    .text-section {
        p {
            color: var(--text-color-2);
            font-size: 0.7rem;
            line-height: 1.25rem;
        }

        h5 {
            color: var(--text-color-1);
        }
    }

    .side-label-title {
        color: var(--text-color-2);
        font-size: 1.05rem;
        line-height: 1.75rem;
        margin-bottom: 0.35em;
        padding-left: 0.35em;
    }

    .route-side-panel {
        display: flex;
        flex-direction: column;
        gap: 0.25em;

        .route-link {
            display: flex;
            align-items: center;
            gap: 0.5em;

            color: #bfbfbf;
            text-decoration: none;
            padding: 0.5em;
            border-radius: 8px;

            &:hover {
                background-color: var(--gray-color-2);
            }

            span {
                font-size: 0.875rem;
                line-height: 1.25rem;
            }
        }

        :global(.route-link-active) {
            background-color: var(--gray-color-2);

            color: #efefef;
        }
    }

    .badge-side-panel {
        position: relative;
        height: fit-content;
        margin-top: auto;

        .badge-wrapper {
            position: relative;
            overflow: hidden;
            width: 100m;
            height: 6em;
            background-color: var(--accent-red);
            border-radius: 8px;
            transition: height 0.3s ease;

            &:hover {
                height: 12em;

                img {
                    top: -20%;
                }
            }

            img {
                position: absolute;
                top: -50%;
                right: -10%;
                width: 120%;
                aspect-ratio: 1 / 1;
                object-fit: cover;

                transition: top 0.3s ease;
            }
        }
    }
</style>
