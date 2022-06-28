<script lang="ts">
    import "../../app.css"
    import { page } from "$app/stores";
    import { WorkerState, WorkerStateToString, WorkerTypeToString, type Job, type PalindromeResult, type WorkerTypeState } from "$lib/job";

    export let jobStatus: Job;

    let createdDate: Date = new Date(jobStatus.created ?? 0);
    let sorterDate: Date = new Date(jobStatus.sortFinish ?? 0);
    let palinDate: Date = new Date(jobStatus.palindromeFinish ?? 0);

    let sorterRuntime: number = (sorterDate.getTime() - createdDate.getTime()) / 1000
    let palinRuntime: number = (palinDate.getTime() - createdDate.getTime()) / 1000

    
    let workerStatus: WorkerTypeState[] = [];
    
    for (const key in jobStatus.workers) {
        if (Object.prototype.hasOwnProperty.call(jobStatus.workers, key)) {
            const element = jobStatus.workers[key];
            workerStatus.push(element);
        }
    }
    
    let workersComplete = 0;
    for (let worker of workerStatus) {
        if (worker.state == 3) workersComplete++
    }
    let percentComplete = (workersComplete/workerStatus.length) * 100
    export let palindromeResult: PalindromeResult | undefined;
    
</script>

{#if jobStatus.error && jobStatus.error != "" }
    We encountered the following error: {jobStatus.error}!
{:else}
    <div>
        <p>Status for job '{jobStatus.id}': {WorkerStateToString(jobStatus.state ?? WorkerState.Failed)}</p>
        <div class="bar">
            <div class="progress" style="width: {percentComplete}%">
                <p class="percentage">{percentComplete.toFixed(0)}%</p>
            </div>
        </div>
        <p>Job information:</p>
        <p>Created: {createdDate.toLocaleString("en-GB")}</p>
        {#if sorterDate.valueOf() > 0}
        <p>Sorter runtime: {sorterRuntime}s</p>
        <p>Sorting finished @ {sorterDate.toLocaleString("en-GB")}</p>
        {/if}
        {#if palinDate.valueOf() > 0}
            <p>Palindrome runtime: {palinRuntime}s</p>
            <p>Finding palindromes finished @ {palinDate.toLocaleString("en-GB")}</p>
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
    .bar {
        margin-top: 10px;
        margin-bottom: 10px;
        height: 20px;
        width: 100%;
        background: rgb(180, 180, 180);
        border-radius: 20px;
    }
    .progress {
        border-radius: 20px;
        display: flex;
        height: 20px;
        background: rgb(181, 255, 172);
        justify-content: center;
    }
    .percentage{
        text-align: center;
        position: absolute;
        left: 50%;
    }
</style>