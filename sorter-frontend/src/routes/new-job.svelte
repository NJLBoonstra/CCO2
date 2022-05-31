<script lang="ts">
    import "../app.css"

    function formSubmit(e: Event) {
        console.log("nothing here yet...")
    }
</script>

<svelte:head>
    <title>New Job</title>
</svelte:head>

<div>
    <p>Select a file to upload, or use the textbox to upload raw text (for debugging)</p>
    <form on:submit|preventDefault={formSubmit}>
        <div>
            <input type="file" multiple>
            <textarea placeholder="Insert text here..."></textarea>
        </div>
        <input type="submit" value="Upload">
    </form>
</div>

<style>
    * {
        text-align: center;
    }
    div {
        display: flex;
        flex-flow: column nowrap;
        align-items: center;
        justify-content: center;
    }
    form>div {
        display: flex;
        flex-flow: row nowrap;
        align-items: stretch;
        justify-content: stretch;
        gap: 20px;
    }
    form>div * {
        flex: 1 1 0;
    }
    form>input[type="submit"] {
        width: 100%;
        padding: 10px;
    }
    textarea {
        max-width: 300px;
        padding: 10px;
        opacity: 50%;
        transition: 350ms;
        font-family: var(--font-family);
        background-color: var(--tertiary-color);
        border: 1px solid black;
        border-radius: var(--border-radius);
        text-align: center;
        transition: var(--transition-time);
        min-height: 200px;
        min-width: 200px;
    }
    textarea:hover, textarea:focus {
        opacity: 100%;
        outline: none;
        background-color: (--pure-white);
    }
</style>