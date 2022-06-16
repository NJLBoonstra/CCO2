<script lang="ts">
    import "../../app.css"
    import { page } from "$app/stores";
    import { jobStateToString, type Job } from "$lib/googlecloud";

    export let jobStatus: Job;
</script>

{#if jobStatus.error }
    We encountered the following error: {jobStatus.error}!
{:else}
    <div>
        <p>Chunk status for job '{$page.params.id}':</p>
        {#each jobStatus.sortState as chunk, id}
            <div>
                <p>Chunk #{id}</p>
                <p>{jobStateToString(chunk)}</p>
            </div>
        {/each}
    </div>
{/if}


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