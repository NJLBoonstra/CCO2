<script lang="ts">
    import "../app.css"

    export const prerender: boolean = true;
</script>

<svelte:head>
    <title>BeanHouse Sorter</title>
</svelte:head>

<h1>Welcome,</h1>

<p>Use the buttons below to either check the status of your job, or to create a new job.</p>

<div>
    <a href="/job-status">Job Status</a>
    <a href="/new-job">New Job</a>
</div>

<style>
    h1 {
        text-align: center;
        color: var(--heading-color);
    }
    p {
        text-align: center;
    }
    div {
        display: flex;
        gap: 50px;
        flex-flow: row nowrap;
        align-items: center;
        justify-content: center;
    }
    a {
        text-align: center;
        padding: 25px;
        min-width: 200px;
        text-decoration: none;
        display: block;
        border-radius: 8px;
        color: var(--accent-color);
        background: linear-gradient(
            90deg,
            var(--primary-color) 0%,
            var(--secondary-color) 50.45%,
            var(--tertiary-color) 61.35%
	    );
        background-size: 400% 400%;
    }
    @keyframes abackground {
        0% {background-position: 0% 50%;}
        50% {background-position: 100% 50%;}
        100% {background-position: 0% 50%;}
    }
    a:hover {
        animation-name: abackground;
        animation-duration: 2s;
        animation-iteration-count: infinite;
        animation-timing-function: ease-in-out;
        animation-direction: alternate;
    }
</style>