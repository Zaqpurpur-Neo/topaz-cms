<script lang="ts">
    import { ChevronRight } from "@lucide/svelte";

    type Crumb = {
        label: string;
        href?: string;
    };

    let { items }: { items: Crumb[] } = $props();
</script>

<nav class="breadcrumb-section" aria-label="breadcrumb">
    <ul class="breadcrumb-list">
        {#each items as item, index}
            <li class="breadcrumb-item">
                <a href={item.href}>{item.label}</a>
                {#if index < items.length - 1}
                    <ChevronRight size={16} />
                {/if}
            </li>
        {/each}
    </ul>
</nav>

<style>
    .breadcrumb-section {
        position: relative;
    }

    .breadcrumb-list {
        display: flex;
        list-style: none;
        padding: 0;
        margin: 0;
        gap: 0.5em;
    }

    .breadcrumb-item {
        display: flex;
        align-items: center;
        gap: 0.5em;
    }

    .breadcrumb-item a {
        text-decoration: none;
        color: #fff;
        font-size: 1.5rem;
        font-weight: 600;

        &:hover {
            text-decoration: underline;
        }
    }
</style>
