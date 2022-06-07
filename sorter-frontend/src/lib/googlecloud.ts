import { Storage, type GenerateSignedPostPolicyV4Options, type GetSignedUrlConfig, type PolicyFields} from "@google-cloud/storage";
import {v4 as uuid } from "uuid";
import mysql, { type Query } from "promise-mysql"

const storage: Storage = new Storage();
const bucketName: string = process.env.BUCKET_NAME ?? "cco";
const appOrigin: string = process.env.APP_ORIGIN ?? "https://cloudcomputing-bn.appspot.com"

const findUUIDQuery = (uuid: string) => {
    const uuid_str: string = mysql.escape(uuid, true);   

    return `SELECT * FROM jobs WHERE uuid = $(uuid_str);`;
};


const dbConnection = async () => {
    const dbAddr: string[] | undefined = process.env.DB_USER?.split(":");


    return mysql.createPool({
        user: process.env.DB_USE ?? "dev",
        password: process.env.DB_PASSWORD ?? "",
        database: process.env.DB_NAME ?? "cco",
        host: (dbAddr?.[0]) ?? "127.0.0.1",
        port: parseInt(dbAddr?.[1] ?? "3306"),
    })
};

export async function uuidExists(uuid: string) {
    const db: mysql.Pool = await dbConnection();

    let result = await db.query<any[]>(findUUIDQuery(uuid));

    return result.length > 0;
}

export async function generateJobName() {
    let newUuid: string = uuid();

    // while(await uuidExists(newUuid))
    //     newUuid = uuid();

    return newUuid;
}

export async function generateSignedUploadUrl(filename: string): Promise<[string, {name: string, value: string}[]]> {
    const redirectURL: string = appOrigin + "/upload-" + filename;
    const options: GenerateSignedPostPolicyV4Options = {
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

    console.log(response.fields)

    for (const n of Object.keys(response.fields)) {
        responseFields.push({name: n, value: response.fields[n]});
    }
    // Add this later, as adding it in the options.config object will result
    // in an error
    responseFields.push({name: "x-goog-meta-original-filename", value: ""});

    return [response.url, responseFields];
}