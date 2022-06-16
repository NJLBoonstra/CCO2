export enum JobState {
    Created,
    Running,
    Completed,
    Failed
};

export type Job = {
	ID?: string,
	State?: JobState,
	sortState?: JobState[],
	palindromeState?: JobState[],
    error?: string,
};

export function jobStateToString(js: JobState) {
    switch (js) {
        case JobState.Created:
            return "Created";
        case JobState.Running:
            return "Running";
        case JobState.Completed:
            return "Completed";
        case JobState.Failed:
            return "Failed";
        default:
            return "???";
    }
}