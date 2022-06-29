import * as gcs from "@google-cloud/storage";
import { v4 as uuid } from "uuid";
import type { Job, PalindromeResult} from "./job";
import { GoogleAuth, IdTokenClient } from "google-auth-library";


const storage: gcs.Storage = new gcs.Storage();
const bucketName: string = process.env.BUCKET_NAME ?? "cco";
const appOrigin: string = process.env.APP_ORIGIN ?? "localhost"
const urlAPI: string = process.env.URL_API ?? "";
const auth = new GoogleAuth();

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

export async function getPalindromeResult(jobID: string): Promise<PalindromeResult> {
    const reqURL: URL = new URL(urlAPI + "/" + jobID + "/palindrome");
    const authReq: IdTokenClient = await auth.getIdTokenClient(urlAPI);

    const response = await authReq.request<PalindromeResult>({url: reqURL.href });
    let data: PalindromeResult;

    if (response.status === 200) {
        data = response.data;
    }
    else
        data = {
            error: "Cloud function returned non-200 code",
        }

    return data
}

export async function getJobList(): Promise<Job[]> {
    const reqURL: URL = new URL(urlAPI + "/all/list")
    const authReq: IdTokenClient = await auth.getIdTokenClient(urlAPI);

    const response = await authReq.request<Job[]>({url: reqURL.href});
    let data: Job[] = [];

    if (response.status === 200) {
        data = response.data;
    }


    return data;

}

export async function getJobStatus(jobID: string): Promise<Job> {
    const reqURL: URL = new URL(urlAPI + "/" + jobID)
    const authReq: IdTokenClient = await auth.getIdTokenClient(urlAPI);

    const response = await authReq.request<Job>({url: reqURL.href});
    let data: Job;
    
    if (response.status === 200) {
        data = response.data;
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

export async function generateSignedDownloadUrl(filename: string): Promise<{filename: string, url: string}> {
    const options: gcs.GetSignedUrlConfig = {
        action: "read",
        expires: 10,
    }
    const [url] = await storage.bucket(bucketName).file(filename).getSignedUrl(options);
    const fname = storage.bucket(bucketName).file(filename).metadata?.["original-filename"] ?? "";

    return {
        filename: fname,
        url: url,
    }
}

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