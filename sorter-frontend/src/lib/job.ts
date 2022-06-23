export enum WorkerState {
    Created,
    Running,
    Reducing,
    Completed,
    Failed
};

export enum WorkerType {
    Sorter,
    Palindrome,
    SorterReduce,
    PalindromeReduce,
}

export type WorkerTypeState = {
    type: WorkerType,
    state: WorkerState,
};

export type Job = {
	id?: string,
	state?: WorkerState,
	workers?: { [id: string]: WorkerTypeState},
    error?: string,
};

export type PalindromeResult = {
    jobId?: string,
    palindromes?: number,
    longestPalindrome?: number,
    error?: string,
}

export function WorkerStateToString(ws: WorkerState): string {
    switch (ws) {
        case WorkerState.Created:
            return "Created";
        case WorkerState.Running:
            return "Running";
        case WorkerState.Completed:
            return "Completed";
        case WorkerState.Reducing:
            return "Reducing";
        case WorkerState.Failed:
            return "Failed";
        default:
            return "???";
    }
}
export function WorkerTypeToString(wt: WorkerType): string {
    switch (wt) {
        case WorkerType.Palindrome:
            return "Palindrome";
        case WorkerType.PalindromeReduce:
            return "Palindrome Reducer";
        case WorkerType.Sorter:
            return "Sort";
        case WorkerType.SorterReduce:
            return "Sort Reducer";
        default:
            return "???";
    }
}