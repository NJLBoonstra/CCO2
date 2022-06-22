<script lang="ts">
    import "../../app.css"
    import { page } from "$app/stores";
    import { WorkerState, WorkerStateToString, WorkerType, WorkerTypeToString, type Job, type PalindromeResult, type WorkerTypeState } from "$lib/job";

    export let jobStatus: Job;
    
    let workerStatus: WorkerTypeState[] = [];

    for (const key in jobStatus.workers) {
        if (Object.prototype.hasOwnProperty.call(jobStatus.workers, key)) {
            const element = jobStatus.workers[key];
            workerStatus.push(element);
        }
    }


    export let palindromeResult: PalindromeResult | undefined;
    
</script>

{#if jobStatus.error && jobStatus.error != "" }
    We encountered the following error: {jobStatus.error}!
{:else}
    <div>
        <p>Status for job '{$page.params.id}': {WorkerStateToString(jobStatus.state ?? WorkerState.Failed)}</p>
        <p>Worker information:</p>
        <table>
            <thead>
                <th>Worker</th>
                <th>Type</th>
                <th>Status</th>
            </thead>
            <tbody>
                {#each workerStatus as s, i}
                <tr>
                    <td>{i}</td>
                    <td>{WorkerTypeToString(s.type)}</td>
                    <td>{WorkerStateToString(s.state)}</td>
                </tr>
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