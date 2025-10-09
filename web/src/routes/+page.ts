import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch }) => {
  const linkReq = await fetch("/api/link");
  const linkJson = await linkReq.json();

  const appsReq = await fetch("/api/apps/", { headers: { Authorization: linkJson.ApiKey } });
  const appsJson = await appsReq.json();

  return {
    domain: linkJson.AccessUrl.split("://")[1],
    apiKey: linkJson.ApiKey,
    apps: appsJson,
  }
};
