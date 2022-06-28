<script lang="ts">
    import "../app.css"

    export const prerender: boolean = true;
</script>

<svelte:head>
    <title>BeanHouse Sorter</title>
</svelte:head>


<div class="main-container">
    <h1>BeanHouse's Sorting as a Service</h1>
    
    <p>Use the buttons below to either check the status of your job, or to create a new job.</p>
    
    <div class="buttons">
        <a href="/job-status">Job Status</a>
        <a href="/new-job">New Job</a>
    </div>
</div>

<style>
    h1 {
        text-align: center;
        color: var(--heading-color);
    }
    p {
        text-align: center;
    }
    .main-container {
        padding: 30px;
    }
    .buttons {
        font-weight: bold;
        display: flex;
        flex-flow: row nowrap;
        align-items: center;
        justify-content: center;
        gap: 50px;
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
        100% {background-position: 100% 50%;}
    }
    a:hover {
        animation-name: abackground;
        animation-fill-mode: forwards;
        animation-duration: 1s;
        animation-direction: normal;
    }
</style>