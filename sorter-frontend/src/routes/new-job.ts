import { generateJobName, generateSignedUploadUrl } from "$lib/googlecloud";
import type { PolicyFields } from "@google-cloud/storage";

export async function get() {
    const jobName: string = await generateJobName();
    let postURL: string;
    let postFields: PolicyFields;
    [postURL, postFields] = await generateSignedUploadUrl(jobName);
   
    return {
      body: { postURL, postFields }
    };
}

// export async function post({params}) {
//     const uuid: string = params.uuid;
//     const fileName: string = params.fileName;

//     // TODO: Link filename and uuid in the db

//     const postURL: string = await generateSignedUploadUrl(uuid);

//     return {
//         body: { postURL }
//     };

// }