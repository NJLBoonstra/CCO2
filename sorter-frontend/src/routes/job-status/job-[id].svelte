<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";


    onMount(async () => {
        // Do some async stuff from endpoints
    });
</script>

Job with ID {$page.params.id} not found...