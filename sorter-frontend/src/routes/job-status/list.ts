import type { RequestHandler } from "@sveltejs/kit";
import { getJobList} from "$lib/googlecloud";
import type { Job } from "$lib/job";


/** @type {import('./__types/items').RequestHandler} */
export async function get() {
    const jobs: Job[] = await getJobList();

    return {
        body: {
            jobs
        }
    };
}