<script lang="ts">
    import "../../app.css"
    import { page } from "$app/stores";
    import { WorkerState, WorkerStateToString, WorkerTypeToString, type Job, type PalindromeResult, type WorkerTypeState } from "$lib/job";

    export let jobStatus: Job;

    let createdDate: Date = new Date(jobStatus.created ?? 0);
    let sorterDate: Date = new Date(jobStatus.sortFinish ?? 0);
    let palinDate: Date = new Date(jobStatus.palindromeFinish ?? 0);

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
        <p>Job information:</p>
        <p>Created: {createdDate.toLocaleString()}</p>
        {#if sorterDate.valueOf() > 0}
            <p>Sorting finished @ {sorterDate.toLocaleString()}</p>
        {/if}
        {#if palinDate.valueOf() > 0}
            <p>Finding palindromes finished @ {palinDate.toLocaleString()}</p>
        {/if}
        {#if palindromeResult && palindromeResult.jobId !== ""}
            <p>The file contains {palindromeResult.palindromes} palindromes and the longest is {palindromeResult.longestPalindrome} characters.</p>
        {/if}
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
    </div>
{/if}


<style>
    div {
        display: flex;
        flex-flow: column wrap;
        max-height: 100%;
    }
    table {
        text-align: left;
    }
    thead {
        text-align: left;
    }
</style>