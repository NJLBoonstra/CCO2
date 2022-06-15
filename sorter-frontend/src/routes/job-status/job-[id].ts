import type { RequestHandler } from "@sveltejs/kit";
import { getJobStatus, type Job } from "$lib/googlecloud";

/** @type {import('./__types/items').RequestHandler} */
export async function get({params}) {
    const jobID: string = params.id;

    const jobStatus: Job = await getJobStatus(jobID);

    return {
        body: {
            jobStatus
        }
    }
}