<script lang="ts">
    import { Gem } from "@lucide/svelte";

    import { routerName, p } from "../router";
    import { isActiveLink } from "sv-router";
</script>

<section class="side-panel">
    <div class="side-panel-header">
        <div class="icon-section">
            <Gem size={18} />
        </div>
        <div class="text-section">
            <h5>Topaz CMS</h5>
            <p>Simple git based CMS</p>
        </div>
    </div>

    <div class="gap-section"></div>

    <div class="route-side-panel">
        {#each routerName as route}
            {#if route.isVisible}
                <a
                    class="route-link"
                    {@attach isActiveLink({ className: "route-link-active" })}
                    href={p(route.path)}
                >
                    <svelte:component this={route.icon} size={16} />
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

<style>
    .side-panel {
        background-color: var(--gray-color-1);
        padding: 0.95em 1em;
        min-width: 16em;
        width: min-content;

        display: flex;
        flex-direction: column;
        height: 100%;
    }

    .gap-section {
        position: relative;
        height: 2em;
    }

    .side-panel-header {
        display: flex;
        align-items: center;

        gap: 0.5em;
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

    .route-side-panel {
        display: flex;
        flex-direction: column;
        gap: 0.5em;

        :global(.route-link-active) {
            background-color: var(--gray-color-2);
        }

        .route-link {
            display: flex;
            align-items: center;
            gap: 0.5em;

            color: var(--text-color);
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
