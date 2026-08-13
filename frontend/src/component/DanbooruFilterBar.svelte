<script lang="ts">
    import { DropdownMenu } from "bits-ui";
    import { Check, ChevronDown, SlidersHorizontal } from "@lucide/svelte";

    interface Props {
        selectedRatings?: string[]; // e.g. ['g', 's']
        limit?: number;
        onChange?: () => void;
    }

    let {
        selectedRatings = $bindable(["g", "s" /* "q", "e" */]),
        limit = $bindable(20),
        onChange = () => {},
    }: Props = $props();

    const ratingMap = [
        { value: "g", label: "General" },
        { value: "s", label: "Sensitive" },
        { value: "q", label: "Questionable" },
        { value: "e", label: "Explicit" },
    ];

    function toggleRating(value: string) {
        if (selectedRatings.includes(value)) {
            if (selectedRatings.length > 1) {
                selectedRatings = selectedRatings.filter((r) => r !== value);
            }
        } else {
            selectedRatings = [...selectedRatings, value];
        }
        onChange();
    }

    function handleLimitInput(e: Event) {
        const val = parseInt((e.target as HTMLInputElement).value, 10);
        if (!isNaN(val) && val >= 1 && val <= 200) {
            limit = val;
            onChange();
        }
    }

    let ratingLabel = $derived.by(() => {
        if (selectedRatings.length === 4) return "All Ratings";
        if (selectedRatings.length === 0) return "None";
        return selectedRatings.map((r) => r.toUpperCase()).join(", ");
    });
</script>

<div class="filter-bar">
    <DropdownMenu.Root>
        <DropdownMenu.Trigger class="filter-dropdown-btn">
            <SlidersHorizontal size={16} />
            <span>Rating: <strong>{ratingLabel}</strong></span>
            <ChevronDown size={14} />
        </DropdownMenu.Trigger>

        <DropdownMenu.Portal>
            <DropdownMenu.Content class="filter-dropdown-menu">
                {#each ratingMap as rating (rating.value)}
                    <DropdownMenu.CheckboxItem
                        checked={selectedRatings.includes(rating.value)}
                        onCheckedChange={() => toggleRating(rating.value)}
                        class="filter-checkbox-item"
                    >
                        <div class="checkbox-indicator">
                            {#if selectedRatings.includes(rating.value)}
                                <Check size={14} />
                            {/if}
                        </div>
                        <span class="rating-name">{rating.label}</span>
                        <span class="rating-code">({rating.value})</span>
                    </DropdownMenu.CheckboxItem>
                {/each}
            </DropdownMenu.Content>
        </DropdownMenu.Portal>
    </DropdownMenu.Root>

    <div class="limit-box">
        <label for="limit-input">Limit:</label>
        <input
            id="limit-input"
            type="number"
            min="1"
            max="200"
            value={limit}
            onchange={handleLimitInput}
        />
    </div>
</div>

<style>
    .filter-bar {
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-bottom: 1rem;
    }

    :global(.filter-dropdown-btn) {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        background-color: #262626;
        border: 1px solid #303030;
        border-radius: 8px;
        padding: 0.5rem 0.85rem;
        color: #efefef;
        font-size: 0.875rem;
        cursor: pointer;
        transition: border-color 0.15s;
    }

    :global(.filter-dropdown-btn:hover) {
        border-color: #3b82f6;
    }

    :global(.filter-dropdown-menu) {
        background-color: #262626;
        border: 1px solid #303030;
        border-radius: 8px;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.4);
        padding: 4px;
        min-width: 180px;
        z-index: 50;
    }

    :global(.filter-checkbox-item) {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        padding: 0.5rem 0.75rem;
        border-radius: 6px;
        color: #efefef;
        font-size: 0.875rem;
        cursor: pointer;
        user-select: none;
    }

    :global(.filter-checkbox-item[data-highlighted]) {
        background-color: #383838;
    }

    .checkbox-indicator {
        width: 18px;
        height: 18px;
        border: 1px solid #545454;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #1e1e1e;
        color: #3b82f6;
    }

    .rating-name {
        flex: 1;
    }

    .rating-code {
        color: #9ca3af;
        font-size: 0.75rem;
        text-transform: uppercase;
    }

    .limit-box {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        background-color: #262626;
        border: 1px solid #303030;
        border-radius: 8px;
        padding: 0.35rem 0.75rem;
        color: #9ca3af;
        font-size: 0.875rem;
    }

    .limit-box input {
        width: 50px;
        background: transparent;
        border: none;
        outline: none;
        color: #efefef;
        font-size: 0.875rem;
        font-weight: 600;
    }

    .limit-box input::-webkit-inner-spin-button,
    .limit-box input::-webkit-outer-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }
</style>
