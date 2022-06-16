<script lang="ts">
    import "../../app.css"
    import { page } from "$app/stores";
    import { JobState, jobStateToString, type Job } from "$lib/googlecloud";

    export let jobStatus: Job;
    let statii: JobState[][];

    if (jobStatus.palindromeState && jobStatus.sortState) {
        statii = jobStatus.sortState.map((v, i) => {
            if (i < jobStatus.palindromeState!.length)
                return [v, jobStatus.palindromeState![i]]
            return [v, JobState.Completed]
        });
    }
</script>

{#if jobStatus.error }
    We encountered the following error: {jobStatus.error}!
{:else}
    <div>
        <p>Chunk status for job '{$page.params.id}':</p>
        {#each statii as chunk, id}
            <div>
                <p>Chunk #{id}</p>
                <p>{jobStateToString(chunk?.[0])}</p>
                <p>{jobStateToString(chunk?.[1])}</p>
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