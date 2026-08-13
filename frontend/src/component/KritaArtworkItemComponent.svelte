<script lang="ts">
    import { p } from "../router";
    let { projectName, artworkName, description, preview, createdAt } =
        $props();

    function formatDate(date: string) {
        // format: May 26, 2025, how can i do this
        const options: Intl.DateTimeFormatOptions = {
            year: "numeric",
            month: "long",
            day: "numeric",
        };
        return new Date(date).toLocaleDateString(undefined, options);
    }
</script>

<a
    class="krita-workspace-item"
    href={p(`/krita-workspace/${projectName}/board/${artworkName}`)}
>
    <div class="wrapper-up">
        <img class="preview" src={preview} alt={projectName} />
    </div>
    <span class="divider"></span>
    <div class="wrapper-down">
        <p class="created-at">{formatDate(createdAt)}</p>
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

        max-width: 24em;
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
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .wrapper-down {
        padding: 0.75em;
    }

    /* ellipsis text if overshot */
    .project-name {
        font-size: 1.25rem;
        font-weight: 600;
        line-height: 1.75;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .description {
        font-size: 0.75rem;
        color: #999;
        line-height: 1.5;
    }

    .created-at {
        font-size: 0.75rem;
        color: #777;
    }
</style>
