<script lang="ts">
    import { goto } from "$app/navigation";
    import "../../app.css"

    export const prerender: boolean = true;

    function formSubmit(e: Event) {
        const formData: FormData = new FormData(e!.target! as HTMLFormElement);

        const jobID: string | undefined = formData.get("jobID")?.toString();

        if (jobID)
            goto("/job-status/job-" + jobID);
    }
</script>

<svelte:head>
    <title>Job Status</title>
</svelte:head>

<div>
    <p>Please enter your job ID:</p>
    <form on:submit|preventDefault={formSubmit}><input name="jobID" type="text"></form>
    <p>And press enter</p>
</div>
<div class=smaller ><a href="/job-status/list">All jobs</a></div>

<style>
    div {
        display: flex;
        flex-flow: column nowrap;
        align-items: center;
    }
    input[type="text"] {
        max-width: 300px;
        padding: 10px;
        opacity: 50%;
        transition: 350ms;
        font-family: var(--font-family);
        background-color: var(--tertiary-color);
        border: 1px solid black;
        border-radius: var(--border-radius);
        text-align: center;
        font-size: larger;
        transition: var(--transition-time);
    }
    p {
        padding: 20px;
    }
    input[type="text"]:hover, input[type="text"]:focus {
        opacity: 100%;
        outline: none;
        background-color: var(--pure-white);
    }
</style>