import * as gcs from "@google-cloud/storage";
import { browser } from "$app/env";
import {stringify, v4 as uuid } from "uuid";
import type { Job } from "./job";


const storage: gcs.Storage = new gcs.Storage();
const bucketName: string = process.env.BUCKET_NAME ?? "cco";
const appOrigin: string = process.env.APP_ORIGIN ?? "https://cloudcomputing-bn.appspot.com"
const urlAPI: string = process.env.URL_API ?? "";

// async function uuidExists(uuid: string) {
//     const obj: gcs.File = storage.bucket(bucketName).file(uuid);

//     return await obj.exists();
// }

export async function generateJobName() {
    let newUuid: string = uuid();

    // while(await uuidExists(newUuid))
    //     newUuid = uuid();

    return newUuid;
}

export async function getJobStatus(jobID: string): Promise<Job> {
    const reqURL: URL = new URL(urlAPI + "/" + jobID)

    const response: Response = await fetch(reqURL.href);
    let data: Job;

    if (response.status === 200) {
        data = await response.json() as Job;
    }
    else
        data = {
            error: "Cloud Function returned non-200 code",
        }


    // try{
    // }
    // catch (e) {
    //     let message: string;
    //     if (e instanceof Error) {
    //         message = e.message;
    //     } else {
    //         message = String(e)
    //     }

    //     data = {
    //         error: message,
    //     };        
    // }
    return data;
}

// export async function generateSignedDownloadUrl(jobID: string) {

// }

export async function generateSignedUploadUrl(filename: string): Promise<[string, {name: string, value: string}[]]> {
    const redirectURL: string = appOrigin + "/upload-" + filename;
    const options: gcs.GenerateSignedPostPolicyV4Options = {
        expires: Date.now() + 15 * 60 * 1000,
        fields: {
            // "x-goog-meta-original-filename": "",
            "success_action_redirect": redirectURL,
        },
        conditions: [
            ["starts-with", "$x-goog-meta-original-filename", ""],
            // {"$success_action_redirect": redirectURL}
        ]
    }

    const [response] = await storage.bucket(bucketName).file(filename)
                    .generateSignedPostPolicyV4(options);

    let responseFields: {name: string, value: string}[] = [];

    for (const n of Object.keys(response.fields)) {
        responseFields.push({name: n, value: response.fields[n]});
    }
    // Add this later, as adding it in the options.config object will result
    // in an error
    responseFields.push({name: "x-goog-meta-original-filename", value: ""});

    return [response.url, responseFields];
}