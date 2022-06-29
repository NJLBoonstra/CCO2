<script lang="ts">
    import "../app.css"

    export let postURL: string;
    export let postFields: {name: string, value:string}[];

    let fileText: string;
    let fileName: string;

    let fileForm: HTMLFormElement;
    let fileInput: HTMLInputElement;
    let textArea: HTMLTextAreaElement;
    // let fileNameElement: HTMLInputElement;

    function resetForm(e: Event) {
        // yeah
        fileInput.value = "";
        fileInput.disabled = false;
        
        textArea.value = "";
        textArea.disabled = false;

        // fileNameElement.value = "";
        // fileNameElement.disabled = false;
    }

    function onChange(e: Event) {
        console.log("onChange called.");
        switch (e.target) {
            case fileInput:
                let v: boolean = true;
                if (fileInput.files?.[0])
                    v = false;

                textArea.disabled = v;
                fileName = fileInput.files?.[0].name ?? "unnamed";

                break;
            case textArea:
                if (textArea.value === "")
                    fileInput.disabled = false;
                else
                    fileInput.disabled = true;
                break;
            default:
                break;
        }
    }


</script>

<svelte:head>
    <title>New Job</title>
</svelte:head>

<div>
    <p>Select a file to upload</p>
    <form method="POST" bind:this={fileForm} action={postURL} enctype="multipart/form-data">
        {#each postFields as {name, value}}
            {#if name === "x-goog-meta-original-filename"}
                <input type="hidden" bind:value={fileName} name={name}>
            {:else}
                <input type="hidden" name={name} value={value}>
            {/if}
        {/each}

        <div>
            <input type="file" name="file" on:change={onChange} bind:this={fileInput} accept=".txt,text/plain">
        </div>
        <div>
            <input type="button" value="Reset" on:click={resetForm} >
            <input class="biggerflex" type="submit" value="Upload">
        </div>
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
    .biggerflex {
        flex: 2 1 0;
    }
</style>