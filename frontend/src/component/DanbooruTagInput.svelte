<script lang="ts">
    import { Combobox } from "bits-ui";
    import { X } from "@lucide/svelte";
    import { tick } from "svelte";
    import { API_URL } from "../common";
    import type { DanbooruAutocompleteItem } from "../types/DanbooruTypes";

    let { tags = $bindable([]), onsearch = () => {} } = $props();

    let searchQuery = $state<string>("");
    let autocompleteResults = $state<DanbooruAutocompleteItem[]>([]);
    let showDropdown = $state<boolean>(false);
    let inputRef = $state<HTMLInputElement | null>(null);
    let debounceTimer: ReturnType<typeof setTimeout> | undefined;

    async function handleInput(event: Event) {
        searchQuery = (event.target as HTMLInputElement).value;
        clearTimeout(debounceTimer);

        if (searchQuery.trim() === "") {
            autocompleteResults = [];
            showDropdown = false;
            return;
        }

        debounceTimer = setTimeout(async () => {
            try {
                const response = await fetch(
                    `${API_URL}/danbooru/autocomplete?q=${encodeURIComponent(searchQuery)}`,
                );
                if (response.ok) {
                    autocompleteResults =
                        (await response.json()) as DanbooruAutocompleteItem[];
                    showDropdown = autocompleteResults.length > 0;
                } else {
                    autocompleteResults = [];
                    showDropdown = false;
                }
            } catch (error) {
                console.error(error);
            }
        }, 300);
    }

    async function addTag(tag: string) {
        const cleanTag = tag.trim().toLowerCase();
        if (cleanTag && !tags.includes(cleanTag)) {
            tags = [...tags, cleanTag];
            onsearch();
        }

        autocompleteResults = [];
        showDropdown = false;

        // Defer input reset until AFTER Bits UI completes its internal item-selection cycle
        await tick();
        setTimeout(() => {
            searchQuery = "";
            if (inputRef) {
                inputRef.value = "";
            }
        }, 0);
    }

    function removeTag(tag: string) {
        tags = tags.filter((t) => t !== tag);
        onsearch();
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === "Enter" && !showDropdown) {
            event.preventDefault();
            if (searchQuery.trim() !== "") {
                addTag(searchQuery.trim());
            } else {
                onsearch();
            }
        } else if (
            event.key === "Backspace" &&
            !searchQuery &&
            tags.length > 0
        ) {
            event.preventDefault();
            removeTag(tags[tags.length - 1]);
        }
    }
</script>

<div class="tags-input-container">
    <div class="tag-box">
        {#each tags as tag (tag)}
            <span class="tag-chip">
                {tag}
                <button
                    class="chip-remove"
                    type="button"
                    aria-label="Remove tag"
                    onclick={() => removeTag(tag)}
                >
                    <X size={12} />
                </button>
            </span>
        {/each}

        <Combobox.Root
            type="single"
            bind:open={showDropdown}
            onValueChange={(val) => {
                if (val) addTag(val);
            }}
        >
            <Combobox.Input
                bind:ref={inputRef}
                oninput={handleInput}
                onkeydown={handleKeyDown}
                placeholder="Add tag..."
                class="tag-input-field"
                onfocus={() => {
                    if (autocompleteResults.length) showDropdown = true;
                }}
            />

            <Combobox.Portal>
                <Combobox.Content class="dropdown">
                    {#each autocompleteResults as item (item.value)}
                        <Combobox.Item
                            value={item.value}
                            label={item.value}
                            class="dropdown-item"
                        >
                            <span class="tag-name">{item.value}</span>
                            {#if item.post_count}
                                <span class="tag-count">
                                    {item.post_count.toLocaleString()} posts
                                </span>
                            {/if}
                        </Combobox.Item>
                    {/each}
                </Combobox.Content>
            </Combobox.Portal>
        </Combobox.Root>
    </div>
</div>

<style>
    .tags-input-container {
        position: relative;
        width: 100%;
    }

    .tag-box {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 0.4rem;
        padding: 0.4rem 0.6rem;
        border: 1px solid #303030;
        border-radius: 8px;
        background: #262626;
        min-height: 46px;
        box-sizing: border-box;
    }

    .tag-box:focus-within {
        border-color: #3b82f6;
    }

    .tag-chip {
        display: inline-flex;
        align-items: center;
        gap: 0.3rem;
        background-color: #545454;
        color: #fff;
        padding: 0.25rem 0.5rem;
        border-radius: 6px;
        font-size: 0.875rem;
        font-weight: 500;
    }

    .chip-remove {
        background: none;
        border: none;
        color: #fff;
        cursor: pointer;
        padding: 0;
        margin: 0;
        display: flex;
        align-items: center;
    }

    .chip-remove:hover {
        color: #ef4444;
    }

    :global(.tag-input-field) {
        border: none !important;
        outline: none !important;
        flex: 1;
        min-width: 140px;
        font-size: 0.95rem;
        padding: 0.25rem 0;
        background: none !important;
        color: #fff !important;
    }

    :global(.dropdown) {
        background: #262626;
        border: 1px solid #303030;
        border-radius: 8px;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
        list-style: none;
        margin-top: 4px;
        padding: 4px 0;
        max-height: 240px;
        overflow-y: auto;
        z-index: 50;
        width: var(--bits-combobox-anchor-width, 100%);
    }

    :global(.dropdown-item) {
        padding: 0.6rem 0.85rem;
        color: #efefef;
        display: flex;
        align-items: center;
        justify-content: space-between;
        cursor: pointer;
        user-select: none;
    }

    :global(.dropdown-item[data-highlighted]) {
        background-color: #383838;
    }

    .tag-name {
        font-size: 0.9rem;
    }

    .tag-count {
        font-size: 0.8rem;
        color: #9ca3af;
    }
</style>
