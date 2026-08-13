<!-- Pagination.svelte -->
<script lang="ts">
    import { ChevronLeft, ChevronRight } from "@lucide/svelte";
    import { Pagination } from "bits-ui";

    interface Props {
        page?: number;
        limit?: number;
        totalCount?: number;
        onpagechange?: (newPage: number) => void;
    }

    let {
        page = $bindable(1),
        limit = 20,
        totalCount = 1000 * 20, // Danbooru cap (~1000 pages max)
        onpagechange = () => {},
    }: Props = $props();
</script>

<div class="pagination-container">
    <Pagination.Root
        count={totalCount}
        perPage={limit}
        bind:page
        onPageChange={(p) => onpagechange(p)}
    >
        {#snippet children({ pages })}
            <div class="pagination-bar">
                <Pagination.PrevButton class="page-nav-btn">
                    <ChevronLeft /> Prev
                </Pagination.PrevButton>

                <div class="pages-list">
                    {#each pages as pageItem (pageItem.key)}
                        {#if pageItem.type === "page"}
                            <Pagination.Page page={pageItem} class="page-btn">
                                {pageItem?.value}
                            </Pagination.Page>
                        {:else}
                            <!-- Standard HTML span for the ellipsis -->
                            <span class="ellipsis">…</span>
                        {/if}
                    {/each}
                </div>

                <Pagination.NextButton class="page-nav-btn">
                    Next <ChevronRight />
                </Pagination.NextButton>
            </div>
        {/snippet}
    </Pagination.Root>
</div>

<style>
    .pagination-container {
        display: flex;
        justify-content: center;
        margin-top: 2rem;
        padding-bottom: 2rem;
    }

    .pagination-bar {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .pages-list {
        display: flex;
        align-items: center;
        gap: 0.25rem;
    }

    :global(.page-nav-btn),
    :global(.page-btn) {
        padding: 0.4rem 0.75rem;
        font-size: 0.875rem;
        background-color: #262626;
        border: 1px solid #303030;
        border-radius: 6px;
        color: #efefef;
        cursor: pointer;
        transition: all 0.15s ease;
    }

    :global(.page-nav-btn:hover:not(:disabled)),
    :global(.page-btn:hover:not([data-selected])) {
        background-color: #383838;
    }

    :global(.page-btn[data-selected]) {
        background-color: #3b82f6 !important;
        border-color: #3b82f6 !important;
        color: #ffffff;
        font-weight: 600;
    }

    :global(.page-nav-btn:disabled) {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .ellipsis {
        color: #6b7280;
        padding: 0 0.4rem;
        font-size: 0.875rem;
        user-select: none;
    }
</style>
