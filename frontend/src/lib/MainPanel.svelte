<script lang="ts">
    import Breadcrumb from "../component/Breadcrumb.svelte";

    const breadcrumbItems = [{ label: "Home", href: "/" }];
</script>

<section class="dashboard-panel">
    <Breadcrumb items={breadcrumbItems} />

    <section class="dashboard-content">
        <h1>Hello There</h1>
    </section>
</section>

<style>
    .dashboard-panel {
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
