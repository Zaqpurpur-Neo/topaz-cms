import type { DanbooruPost } from "./types/DanbooruTypes";

export function getMediumPreviewUrl(post: DanbooruPost): string | null {
  if (post.is_banned || post.is_deleted) return null;

  if (post.media_asset?.variants) {
    const variant360 = post.media_asset.variants.find(
      (v) => v.type === "360x360",
    );
    if (variant360?.url) return variant360.url;

    const variant720 = post.media_asset.variants.find(
      (v) => v.type === "720x720",
    );
    if (variant720?.url) return variant720.url;

    const variant180 = post.media_asset.variants.find(
      (v) => v.type === "180x180",
    );
    if (variant180?.url) return variant180.url;
  }

  return post.large_file_url ?? post.preview_file_url ?? post.file_url ?? null;
}
