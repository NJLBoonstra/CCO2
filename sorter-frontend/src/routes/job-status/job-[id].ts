import type { RequestHandler } from "@sveltejs/kit";
import { generateSignedDownloadUrl, getJobStatus, getPalindromeResult} from "$lib/googlecloud";
import { WorkerState, type Job, type PalindromeResult } from "$lib/job";


/** @type {import('./__types/items').RequestHandler} */
export async function get({params}) {
    const jobID: string = params.id;
    const jobStatus: Job = await getJobStatus(jobID);

    let palindromeResult: PalindromeResult = await getPalindromeResult(jobID);
    let filename: string = "";
    let url: string = "";

    if (jobStatus.state === WorkerState.Completed) {
        let r = await generateSignedDownloadUrl(jobID);
        filename = r.filename;
        url = r.url;
    }
    // let 

    return {
        body: {
            jobStatus, palindromeResult, filename, url
        }
    }
}