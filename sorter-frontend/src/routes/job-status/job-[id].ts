import type { RequestHandler } from "@sveltejs/kit";

/** @type {import('./__types/items').RequestHandler} */
export async function get() {
    const chunks = [true, false, true, false, true, false];
   
    return {
      body: { chunks }
    };
  }