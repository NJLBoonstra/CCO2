import type { RequestHandler } from "@sveltejs/kit";
import { getJobStatus, getPalindromeResult} from "$lib/googlecloud";
import type { Job, PalindromeResult } from "$lib/job";


/** @type {import('./__types/items').RequestHandler} */
export async function get({params}) {
    const jobID: string = params.id;
    const jobStatus: Job = await getJobStatus(jobID);

    let palindromeResult: PalindromeResult = await getPalindromeResult(jobID);


    // let 

    return {
        body: {
            jobStatus, palindromeResult
        }
    }
}