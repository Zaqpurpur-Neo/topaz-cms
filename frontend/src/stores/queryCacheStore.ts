// queryCacheStore.ts
import { get, writable } from "svelte/store";
import type { DanbooruPost } from "../types/DanbooruTypes";

export interface CacheEntry {
  posts: DanbooruPost[];
  timestamp: number;
}

// Key = query string (e.g. "tags=miku&page=1&limit=20")
// Value = { posts, timestamp }
export const queryCache = writable<Record<string, CacheEntry>>({});

// Global setting to auto-fetch or rely on cache
export const autoFetchSetting = writable<boolean>(false);

export const danbooruQueryString = writable<string>("");

// Synchronous getter for router.ts
export function getDanbooruQuery(): string {
  return get(danbooruQueryString);
}
