<!-- PostCard.svelte -->
<script lang="ts">
    import type { DanbooruPost } from "../types/DanbooruTypes";

    let {
        post,
        previewUrl,
        onclick,
    }: {
        post: DanbooruPost;
        previewUrl: string | null;
        onclick?: () => void;
    } = $props();
</script>

<div
    class="post-card"
    role="button"
    tabindex="0"
    {onclick}
    onkeydown={(e) => e.key === "Enter" && onclick?.()}
>
    {#if post.is_banned}
        <div class="placeholder banned">Banned</div>
    {:else if post.is_deleted}
        <div class="placeholder deleted">Deleted</div>
    {:else if previewUrl}
        <div class="img-wrapper">
            <img src={previewUrl} alt="Danbooru #{post.id}" loading="lazy" />
        </div>
    {:else}
        <div class="placeholder restricted">Restricted / Gold Only</div>
    {/if}

    <div class="post-meta">
        <span>#{post.id}</span>
        <span class="rating-tag rating-{post.rating}">{post.rating}</span>
    </div>
</div>

<style>
    .post-card {
        background: #262626;
        border: 1px solid #303030;
        border-radius: 8px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
        cursor: pointer;
        transition:
            transform 0.15s ease,
            border-color 0.15s ease;

        &:hover {
            transform: translateY(-2px);
            border-color: #404040;
        }
    }

    .img-wrapper {
        width: 100%;
        height: 220px;
        background: #1a1a1a;
        overflow: hidden;

        img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
        }
    }

    .placeholder {
        height: 220px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 0.85rem;
        font-weight: 600;

        &.banned {
            background-color: #3f1212;
            color: #fca5a5;
        }
        &.deleted {
            background-color: #27272a;
            color: #a1a1aa;
        }
        &.restricted {
            background-color: #312e81;
            color: #c7d2fe;
        }
    }

    .post-meta {
        padding: 0.5rem 0.75rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 0.8rem;
        background: #262626;
    }

    .rating-tag {
        text-transform: uppercase;
        font-weight: 700;
        font-size: 0.7rem;
        padding: 2px 6px;
        border-radius: 4px;

        &.rating-g,
        &.rating-s {
            background: #166534;
            color: #dcfce7;
        }
        &.rating-q {
            background: #854d0e;
            color: #fef9c3;
        }
        &.rating-e {
            background: #991b1b;
            color: #fee2e2;
        }
    }
</style>
