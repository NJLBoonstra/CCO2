<script lang="ts">
    import "../app.css"
    import { page } from "$app/stores";
    export const prerender: boolean = true;
</script>

Succesfully uploaded file!
Job ID: <a href='/job-status/job-{$page.params.id}'>{$page.params.id}</a>
