<script lang="ts">
    import { FolderOpenDot } from "@lucide/svelte";

    import { API_URL } from "../common";
    import { p } from "../router";
    let { projectName, artworkName, description, preview } = $props();
</script>

<a
    class="krita-workspace-item"
    href={p(`/krita-workspace/${projectName}/artwork/${artworkName}`)}
>
    <div class="wrapper-up">
        <img class="preview" src={preview} alt={projectName} />
    </div>
    <span class="divider"></span>
    <div class="wrapper-down">
        <p class="project-name">{artworkName.replaceAll("-", " ")}</p>
        <p class="description">{description}</p>
    </div>
</a>

<style>
    .krita-workspace-item {
        position: relative;
        display: flex;
        flex-direction: column;
        background-color: var(--gray-color-1);
        border: 1px solid #262626;
        border-radius: 8px;
        width: 100%;
        text-decoration: none;
        color: #fff;

        max-width: 20em;
        overflow: hidden;
    }

    .divider {
        width: 100%;
        height: 1px;
        background-color: #262626;
    }

    .wrapper-up {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
        height: 14em;
    }

    .wrapper-up img {
        width: 80%;
        height: 100%;
        object-fit: cover;
    }

    .wrapper-down {
        padding: 1em;
    }

    .description {
        font-size: 0.75rem;
        color: #bbb;
        line-height: 1.5;
    }
</style>
