<script lang="ts">
    import { WorkerState, WorkerStateToString, WorkerType, WorkerTypeToString } from "$lib/job";


    export let state: WorkerState;
    export let type: WorkerType;
    export let idx: number;

    
</script>

<tr>
    <td>{idx}</td>
    <td>{WorkerTypeToString(type)}</td>
    <td>{WorkerStateToString(state)}</td>
</tr>
