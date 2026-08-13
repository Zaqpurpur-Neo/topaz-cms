// danbooruStore.ts (Classic Store)
import { writable } from "svelte/store";
import type { DanbooruPost } from "../types/DanbooruTypes";

export interface DanbooruCacheState {
  posts: DanbooruPost[];
  lastQueryKey: string;
  hasLoadedOnce: boolean;
  autoFetchOnTabOpen: boolean;
}

export const danbooruCache = writable<DanbooruCacheState>({
  posts: [],
  lastQueryKey: "",
  hasLoadedOnce: false,
  autoFetchOnTabOpen: false,
});
