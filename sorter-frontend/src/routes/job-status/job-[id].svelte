<script lang="ts">
    import { page } from "$app/stores";
    export let chunks: boolean[];
    export let error: string;

    console.log(chunks);
    console.log(error);
</script>

{#if error}
    We encountered the following error: {error}!
{/if}

<div>
    {#if chunks?.length < 1}
        <p>Job with ID {$page.params.id} not found...</p>
    {:else}
        <p>Chunk status for job '{$page.params.id}':</p>
        {#each chunks as chunk, id}
            <div>
                <p>Chunk #{id}</p>
                {#if chunk}
                <p>Done</p>
                {:else}
                <p>Waiting...<p>
                {/if}
            </div>
        {/each}
    {/if}
</div>

<style>
    div {
        display: flex;
        flex-flow: column wrap;
    }
    div>div {
        display: flex;
        flex-flow: row nowrap;
        gap: 15px;
        justify-content: stretch;
    }
</style>