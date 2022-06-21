<script lang="ts">
    import "../../app.css"
    import { page } from "$app/stores";
    import { WorkerState, WorkerType, type Job, type PalindromeResult, type WorkerTypeState } from "$lib/job";
import Workerstateelement from "$lib/components/workerstateelement.svelte";

    export let jobStatus: Job;
    console.log(jobStatus);
    
    let workerStatus: WorkerTypeState[] = [];

    for (const key in jobStatus.workers) {
        if (Object.prototype.hasOwnProperty.call(jobStatus.workers, key)) {
            const element = jobStatus.workers[key];
            console.log(element);
            workerStatus.push(element);
        }
    }


    export let palindromeResult: PalindromeResult | undefined;
    
</script>

{#if jobStatus.error && jobStatus.error != "" }
    We encountered the following error: {jobStatus.error}!
{:else}
    <div>
        <p>Status for job '{$page.params.id}': {jobStatus.state ?? WorkerState.Failed}</p>
        <p>Worker information:</p>
        <table>
            <thead>
                <th>Worker</th>
                <th>Type</th>
                <th>Status</th>
            </thead>
            <tbody>
                {#each workerStatus as s, i}
                    <Workerstateelement idx={i} state={s.State} type={s.Type}></Workerstateelement>
                {/each}
            </tbody>
        </table>
        {#if palindromeResult && palindromeResult.jobId !== ""}
            <p>The file contains {palindromeResult.palindromes} palindromes and the longest is {palindromeResult.longestPalindrome} characters.</p>
        {/if}
    </div>
{/if}


<style>
    div {
        display: flex;
        flex-flow: column wrap;
    }
    table {
        text-align: left;
    }
    thead {
        text-align: left;
    }
</style>