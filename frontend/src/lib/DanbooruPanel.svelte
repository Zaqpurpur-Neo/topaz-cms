<script lang="ts">
    import { onMount } from "svelte";
    import { API_URL } from "../common";

    import { get } from "svelte/store";
    import {
        queryCache,
        autoFetchSetting,
        danbooruQueryString,
    } from "../stores/queryCacheStore";

    import Breadcrumb from "../component/Breadcrumb.svelte";
    import DanbooruFilterBar from "../component/DanbooruFilterBar.svelte";
    import DanbooruTagInput from "../component/DanbooruTagInput.svelte";
    import DanbooruPagination from "../component/DanbooruPagination.svelte";
    import DanbooruPostCard from "../component/DanbooruPostCard.svelte";

    import type { DanbooruPost } from "../types/DanbooruTypes";
    import DanbooruPostModal from "../component/DanbooruPostModal.svelte";
    import { getMediumPreviewUrl } from "../helper";

    type Rating = "g" | "s" | "q" | "e";
    const DEFAULT_RATINGS: Rating[] = ["g" /* "s", "q", "e"*/];
    const DEFAULT_LIMIT = 20;
    const DEFAULT_PAGE = 1;

    let selectedTags = $state<string[]>([]);
    let selectedRatings = $state<(Rating | string)[]>([...DEFAULT_RATINGS]);
    let limit = $state<number>(DEFAULT_LIMIT);
    let currentPage = $state<number>(DEFAULT_PAGE);
    let posts = $state<DanbooruPost[]>([]);
    let selectedPost = $state<DanbooruPost | null>(null);

    let isLoading = $state<boolean>(false);
    let isHydrated = $state<boolean>(false);

    function hydrateStateFromUrl() {
        danbooruQueryString.set(window.location.search);
        const params = new URLSearchParams(window.location.search);

        const tagsParam = params.get("tags");
        if (tagsParam) {
            selectedTags = tagsParam
                .trim()
                .split(/[ +]/) // Split on literal '+' or space delimiter
                .filter(Boolean)
                .map((tag) => {
                    try {
                        return decodeURIComponent(tag);
                    } catch {
                        return tag; // Fallback if string is malformed
                    }
                });
        } else {
            selectedTags = [];
        }

        const ratingsParam = params.get("ratings");
        if (ratingsParam !== null) {
            const parsed = ratingsParam
                .split(",")
                .filter((r) => ["g", "s", "q", "e"].includes(r));
            selectedRatings = parsed.length > 0 ? parsed : [...DEFAULT_RATINGS];
        } else {
            selectedRatings = [...DEFAULT_RATINGS];
        }

        const pageParam = params.get("page");
        if (pageParam !== null) {
            const parsedPage = parseInt(pageParam, 10);
            if (!isNaN(parsedPage) && parsedPage >= 1) currentPage = parsedPage;
        }

        const limitParam = params.get("limit");
        if (limitParam !== null) {
            const parsedLimit = parseInt(limitParam, 10);
            if (!isNaN(parsedLimit) && parsedLimit >= 1 && parsedLimit <= 200)
                limit = parsedLimit;
        }
    }

    function updateUrlFromState(replaceHistory = false) {
        if (!isHydrated) return;

        const params = new URLSearchParams();

        if (selectedTags.length > 0) {
            params.set(
                "tags",
                selectedTags.map((t) => encodeURIComponent(t)).join("+"),
            );
        }

        const currentRatingStr = [...selectedRatings].sort().join(",");
        const defaultRatingStr = [...DEFAULT_RATINGS].sort().join(",");

        if (currentRatingStr !== defaultRatingStr) {
            params.set("ratings", selectedRatings.join(","));
        }

        if (currentPage > 1) {
            params.set("page", currentPage.toString());
        }

        if (limit !== 20) {
            params.set("limit", limit.toString());
        }

        const queryString = params.toString();
        const fullQueryStr = queryString ? `?${queryString}` : "";

        danbooruQueryString.set(fullQueryStr);

        const currentPath = window.location.pathname;

        const newUrl = queryString
            ? `${currentPath}?${queryString}`
            : currentPath;

        if (replaceHistory) {
            window.history.replaceState(null, "", newUrl);
        } else {
            window.history.pushState(null, "", newUrl);
        }
    }

    let currentCacheKey = $derived.by(() => {
        const query = [...selectedTags];

        if (selectedRatings.length > 0 && selectedRatings.length < 4) {
            query.push(`rating:${selectedRatings.join(",")}`);
        }

        const tagPart = query.map((tag) => encodeURIComponent(tag)).join("+");
        return `tags=${tagPart}&page=${currentPage}&limit=${limit}`;
    });

    function clearQueryCache() {
        queryCache.set({});
        loadPosts(true); // Force fetch current view
    }

    async function loadPosts(forceRefresh = false) {
        const key = currentCacheKey;
        const cache = get(queryCache);
        const shouldAutoFetch = get(autoFetchSetting);

        if (!forceRefresh && !shouldAutoFetch && cache[key]) {
            posts = cache[key].posts;
            return;
        }

        isLoading = true;
        try {
            const url = `${API_URL}/danbooru/posts?${key}`;
            const res = await fetch(url);

            if (res.ok) {
                const freshPosts = (await res.json()) as DanbooruPost[];
                posts = freshPosts;

                queryCache.update((prev) => ({
                    ...prev,
                    [key]: {
                        posts: freshPosts,
                        timestamp: Date.now(),
                    },
                }));
            }
        } catch (err) {
            console.error("Fetch failed:", err);
        } finally {
            isLoading = false;
        }
    }

    function handleFilterChange() {
        currentPage = 1;
        updateUrlFromState();
        loadPosts();
    }

    function handlePageChange(page: number) {
        currentPage = page;
        updateUrlFromState();

        clearQueryCache();
        window.scrollTo({ top: 0, behavior: "smooth" });
    }

    onMount(() => {
        hydrateStateFromUrl();
        isHydrated = true;

        loadPosts();

        const handlePopState = () => {
            hydrateStateFromUrl();
            loadPosts();
        };

        window.addEventListener("popstate", handlePopState);

        return () => {
            window.removeEventListener("popstate", handlePopState);
        };
    });

    const breadcrumbItems = $derived([
        {
            label: "Home",
            href: "/",
        },
        {
            label: "Danbooru",
            href: "/danbooru",
        },
    ]);

    function openPostModal(post: DanbooruPost) {
        selectedPost = post;
    }
</script>

<section class="danbooru-panel">
    <Breadcrumb items={breadcrumbItems} />

    <div class="danbooru-content">
        <div class="control-section">
            <DanbooruTagInput
                bind:tags={selectedTags}
                onsearch={handleFilterChange}
            />

            <DanbooruFilterBar
                bind:selectedRatings
                bind:limit
                onChange={handleFilterChange}
            />
        </div>

        <div class="danbooru-results">
            {#if isLoading}
                <p class="status-msg">Loading posts...</p>
            {:else if posts.length === 0}
                <p class="status-msg">No posts found.</p>
            {:else}
                <div class="posts-grid">
                    {#each posts as post (post.id)}
                        <DanbooruPostCard
                            {post}
                            previewUrl={getMediumPreviewUrl(post)}
                            onclick={() => openPostModal(post)}
                        />
                    {/each}
                </div>

                <DanbooruPagination
                    bind:page={currentPage}
                    {limit}
                    onpagechange={handlePageChange}
                />
            {/if}
        </div>
    </div>

    <DanbooruPostModal bind:post={selectedPost} />
</section>

<style>
    .danbooru-panel {
        position: relative;
        width: 100%;
        height: 100dvh;
        padding: 1.75em;
        box-sizing: border-box;
        padding-bottom: 0;

        display: flex;
        flex-direction: column;
        gap: 1.7em;
        padding-bottom: 0;
        overflow-y: auto;
    }

    .control-section {
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        margin-bottom: 1.5rem;
    }

    .danbooru-results {
        position: relative;
        width: 100%;
        height: 100%;
        overflow-y: auto;
    }

    .posts-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
        gap: 1rem;
    }

    .status-msg {
        text-align: center;
        color: #9ca3af;
        margin: 3rem 0;
    }
</style>
