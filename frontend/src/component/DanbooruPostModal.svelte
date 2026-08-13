<!-- PostModal.svelte -->
<script lang="ts">
    import { Dialog } from "bits-ui";
    import {
        X,
        ExternalLink,
        Heart,
        Tag,
        FileText,
        Image as ImageIcon,
    } from "@lucide/svelte";
    import type { DanbooruPost } from "../types/DanbooruTypes";

    let {
        post = $bindable<DanbooruPost | null>(),
    }: {
        post: DanbooruPost | null;
    } = $props();

    // Dialog open state tied to whether post is selected
    let isOpen = $derived(post !== null);

    function handleOpenChange(open: boolean) {
        if (!open) post = null;
    }

    // Get best high-res image for modal (large_file_url or file_url)
    let modalImageUrl = $derived.by(() => {
        if (!post || post.is_banned || post.is_deleted) return null;
        return (
            post.large_file_url ??
            post.file_url ??
            post.preview_file_url ??
            null
        );
    });

    // Extract tag groups
    let artistTags = $derived(
        post?.tag_string_artist?.split(" ").filter(Boolean) ?? [],
    );
    let characterTags = $derived(
        post?.tag_string_character?.split(" ").filter(Boolean) ?? [],
    );
    let copyrightTags = $derived(
        post?.tag_string_copyright?.split(" ").filter(Boolean) ?? [],
    );
    let generalTags = $derived(
        post?.tag_string_general?.split(" ").filter(Boolean) ?? [],
    );
</script>

<Dialog.Root open={isOpen} onOpenChange={handleOpenChange}>
    <Dialog.Portal>
        <!-- Dark Backdrop -->
        <Dialog.Overlay class="modal-overlay" />

        <!-- Modal Content -->
        <Dialog.Content class="modal-content">
            {#if post}
                <!-- Left: Image / Media Container -->
                <div class="media-container">
                    {#if post.is_banned}
                        <div class="modal-placeholder banned">
                            <ImageIcon size={48} />
                            <p>
                                This post is <strong>Banned</strong> on Danbooru.
                            </p>
                        </div>
                    {:else if post.is_deleted}
                        <div class="modal-placeholder deleted">
                            <ImageIcon size={48} />
                            <p>
                                This post was <strong>Deleted</strong> from Danbooru.
                            </p>
                        </div>
                    {:else if modalImageUrl}
                        <img src={modalImageUrl} alt="Danbooru #{post.id}" />
                    {:else}
                        <div class="modal-placeholder restricted">
                            <ImageIcon size={48} />
                            <p>
                                Gold / Paid account required to view original
                                file.
                            </p>
                        </div>
                    {/if}
                </div>

                <!-- Right: Instagram-Style Sidebar -->
                <div class="details-container">
                    <!-- Top Header -->
                    <div class="details-header">
                        <div class="post-title">
                            <Dialog.Title class="title-text"
                                >Post #{post.id}</Dialog.Title
                            >
                            <span class="rating-badge rating-{post.rating}">
                                {post.rating?.toUpperCase()}
                            </span>
                        </div>

                        <Dialog.Close
                            class="close-btn"
                            aria-label="Close dialog"
                        >
                            <X size={18} />
                        </Dialog.Close>
                    </div>

                    <!-- Middle: Scrollable Metadata & Tags -->
                    <div class="details-body">
                        <!-- Stats Row -->
                        <div class="stats-row">
                            <div class="stat-item" title="Score">
                                <Heart size={14} class="stat-icon" />
                                <span>{post.score ?? 0}</span>
                            </div>
                            <div class="stat-item" title="Resolution">
                                <ImageIcon size={14} class="stat-icon" />
                                <span
                                    >{post.image_width ?? 0} × {post.image_height ??
                                        0}</span
                                >
                            </div>
                            <div class="stat-item" title="File Extension">
                                <FileText size={14} class="stat-icon" />
                                <span
                                    >{post.file_ext?.toUpperCase() ??
                                        "UNKNOWN"}</span
                                >
                            </div>
                        </div>

                        <hr class="divider" />

                        <!-- Artist Tags -->
                        {#if artistTags.length > 0}
                            <div class="tag-section">
                                <span class="tag-group-title artist"
                                    >Artist</span
                                >
                                <div class="tag-chips">
                                    {#each artistTags as tag}
                                        <span class="chip artist-chip"
                                            >{tag}</span
                                        >
                                    {/each}
                                </div>
                            </div>
                        {/if}

                        <!-- Character Tags -->
                        {#if characterTags.length > 0}
                            <div class="tag-section">
                                <span class="tag-group-title character"
                                    >Character</span
                                >
                                <div class="tag-chips">
                                    {#each characterTags as tag}
                                        <span class="chip character-chip"
                                            >{tag}</span
                                        >
                                    {/each}
                                </div>
                            </div>
                        {/if}

                        <!-- Copyright / Series Tags -->
                        {#if copyrightTags.length > 0}
                            <div class="tag-section">
                                <span class="tag-group-title copyright"
                                    >Copyright</span
                                >
                                <div class="tag-chips">
                                    {#each copyrightTags as tag}
                                        <span class="chip copyright-chip"
                                            >{tag}</span
                                        >
                                    {/each}
                                </div>
                            </div>
                        {/if}

                        <!-- General Tags -->
                        {#if generalTags.length > 0}
                            <div class="tag-section">
                                <span class="tag-group-title general"
                                    >General Tags</span
                                >
                                <div class="tag-chips">
                                    {#each generalTags.slice(0, 30) as tag}
                                        <span class="chip general-chip"
                                            >{tag}</span
                                        >
                                    {/each}
                                    {#if generalTags.length > 30}
                                        <span class="chip more-chip"
                                            >+{generalTags.length - 30} more</span
                                        >
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>

                    <!-- Bottom Footer -->
                    <div class="details-footer">
                        <a
                            href="https://danbooru.donmai.us/posts/{post.id}"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="danbooru-link"
                        >
                            <span>View on Danbooru</span>
                            <ExternalLink size={14} />
                        </a>
                    </div>
                </div>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<style>
    :global(.modal-overlay) {
        position: fixed;
        inset: 0;
        z-index: 999;
        background-color: rgba(0, 0, 0, 0.85);
        backdrop-filter: blur(6px);
    }

    :global(.modal-content) {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        z-index: 1000;
        width: 90vw;
        max-width: 1100px;
        height: 85vh;
        max-height: 750px;
        background-color: #121212;
        border: 1px solid #27272a;
        border-radius: 12px;
        display: flex;
        overflow: hidden;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.7);
        color: #efefef;
        font-family:
            system-ui,
            -apple-system,
            sans-serif;
    }

    /* Left Image Container */
    .media-container {
        flex: 1;
        background-color: #050505;
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        position: relative;
    }

    .media-container img {
        max-width: 100%;
        max-height: 100%;
        object-fit: contain;
    }

    .modal-placeholder {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
        color: #a1a1aa;
        text-align: center;
        padding: 2rem;
    }

    .modal-placeholder.banned {
        color: #fca5a5;
    }
    .modal-placeholder.deleted {
        color: #a1a1aa;
    }
    .modal-placeholder.restricted {
        color: #c7d2fe;
    }

    /* Right Details Sidebar */
    .details-container {
        width: 360px;
        min-width: 360px;
        background-color: #18181b;
        border-left: 1px solid #27272a;
        display: flex;
        flex-direction: column;
    }

    .details-header {
        padding: 1rem;
        border-bottom: 1px solid #27272a;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .post-title {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    :global(.title-text) {
        font-size: 1.1rem;
        font-weight: 700;
        margin: 0;
    }

    .rating-badge {
        font-size: 0.65rem;
        font-weight: 700;
        padding: 2px 6px;
        border-radius: 4px;
    }

    .rating-G,
    .rating-S {
        background: #166534;
        color: #dcfce7;
    }
    .rating-Q {
        background: #854d0e;
        color: #fef9c3;
    }
    .rating-E {
        background: #991b1b;
        color: #fee2e2;
    }

    :global(.close-btn) {
        background: #27272a;
        border: none;
        color: #a1a1aa;
        padding: 6px;
        border-radius: 6px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background 0.15s;
    }

    :global(.close-btn:hover) {
        background: #3f3f46;
        color: #fff;
    }

    /* Body */
    .details-body {
        flex: 1;
        padding: 1rem;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
    }

    .stats-row {
        display: flex;
        justify-content: space-between;
        font-size: 0.8rem;
        color: #a1a1aa;
    }

    .stat-item {
        display: flex;
        align-items: center;
        gap: 0.35rem;
    }

    .divider {
        border: 0;
        border-top: 1px solid #27272a;
        margin: 0.25rem 0;
    }

    .tag-section {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }

    .tag-group-title {
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    .tag-group-title.artist {
        color: #fca5a5;
    }
    .tag-group-title.character {
        color: #86efac;
    }
    .tag-group-title.copyright {
        color: #d8b4fe;
    }
    .tag-group-title.general {
        color: #93c5fd;
    }

    .tag-chips {
        display: flex;
        flex-wrap: wrap;
        gap: 0.25rem;
    }

    .chip {
        font-size: 0.75rem;
        padding: 2px 6px;
        border-radius: 4px;
        background: #27272a;
        color: #d4d4d8;
    }

    .artist-chip {
        background: #451a1a;
        color: #fca5a5;
    }
    .character-chip {
        background: #14532d;
        color: #bbf7d0;
    }
    .copyright-chip {
        background: #3b0764;
        color: #e9d5ff;
    }
    .general-chip {
        background: #1e293b;
        color: #cbd5e1;
    }
    .more-chip {
        background: #18181b;
        color: #71717a;
        border: 1px solid #27272a;
    }

    /* Footer */
    .details-footer {
        padding: 1rem;
        border-top: 1px solid #27272a;
    }

    .danbooru-link {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        width: 100%;
        padding: 0.6rem;
        background-color: #27272a;
        border-radius: 6px;
        color: #efefef;
        text-decoration: none;
        font-size: 0.85rem;
        font-weight: 600;
        transition: background 0.15s;
    }

    .danbooru-link:hover {
        background-color: #3b82f6;
    }

    /* Responsive Mobile View */
    @media (max-width: 768px) {
        :global(.modal-content) {
            flex-direction: column;
            height: 90vh;
        }

        .details-container {
            width: 100%;
            min-width: 100%;
            height: 50%;
            border-left: none;
            border-top: 1px solid #27272a;
        }
    }
</style>
