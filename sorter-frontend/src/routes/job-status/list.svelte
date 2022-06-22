<script lang="ts">
import { type Job, WorkerStateToString } from "$lib/job";
    import "../../app.css"

    export let jobs: Job[];
</script>

<table>
    <thead>
        <tr>
            <th>Job UUID</th>
            <th>Overall Job status</th>
            <th># Workers</th>
        </tr>
    </thead>
    <tbody>
        {#each jobs as job}
        <tr>
                    
                    <td><a href={"/job-status/job-" + job.id}>{job.id}</a></td>
                    {#if job.state !== undefined}
                        <td>{WorkerStateToString(job.state)}</td>
                    {:else}
                        <td>No information</td>  
                    {/if}
                    {#if job.workers}
                        <td>{Object.keys(job.workers).length}</td>
                    {:else}
                        <td>0</td>
                    {/if}
                </tr>
           
        {/each}
    </tbody>
</table>

<style>
    table {
        margin: 0 auto;
        border-collapse: collapse;
    }
    thead {
        margin: 0 0 5px 0;
        text-align: left;
        border-bottom: 1px solid black;
    }
    td {
        text-align: center;
    }
    th, td {
        padding: 5px;
    }
    tbody tr:hover {
        background-color: rgba(255, 62, 0, 0.4);
    }
</style>