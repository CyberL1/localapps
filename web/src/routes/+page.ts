import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
  const dataReq = await fetch("/api/apps/");
  return { apps: await dataReq.json() }
};
