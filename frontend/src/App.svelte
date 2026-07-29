<script lang="ts">
    import SidePanel from "./lib/SidePanel.svelte";
    import { Router } from "sv-router";

    import "./router";
</script>

<main>
    <SidePanel />

    <Router />
</main>

<style>
    main {
        display: flex;
        width: 100%;
        height: 100%;
    }
</style>
