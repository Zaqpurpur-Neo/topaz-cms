<script lang="ts">
    import Breadcrumb from "../component/Breadcrumb.svelte";

    const breadcrumbItems = [
        {
            label: "Home",
            href: "/",
        },
        {
            label: "Media",
            href: "/media",
        },
    ];
</script>

<section class="media-panel">
    <Breadcrumb items={breadcrumbItems} />
</section>

<style>
    .media-panel {
        position: relative;
        width: 100%;
        height: 100%;
        padding: 1.75em;
        box-sizing: border-box;

        display: flex;
        flex-direction: column;
        gap: 1.7em;
        padding-bottom: 0;
    }
</style>
