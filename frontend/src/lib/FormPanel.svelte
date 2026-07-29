<script>
    import { API_URL } from "../common";

    async function handleSubmit(event) {
        const form = event.target;
        const formData = new FormData(form);
        const response = await fetch(API_URL, {
            method: "POST",
            body: formData,
        });
    }
</script>

<section class="create-post-section">
    <form method="POST" on:submit|preventDefault>
        <input type="text" name="title" placeholder="Title" />
        <textarea placeholder="Description" name="description"></textarea>
        <input type="text" name="slug_name" placeholder="SlugName" />
        <input type="text" name="instagram_url" placeholder="Instagram URL" />
        <input type="date" name="date" placeholder="Date" />

        <button type="submit">Submit</button>
    </form>
</section>

<style>
    /* Styles here */
</style>
