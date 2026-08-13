// types.ts

/** Danbooru Tag Categories */
export enum DanbooruTagCategory {
  General = 0,
  Artist = 1,
  Copyright = 3,
  Character = 4,
  Meta = 5,
}

export interface DanbooruAutocompleteItem {
  type: "tag" | "metatag" | "user" | "pool" | "wiki" | "forum_topic";
  label: string; // Formatted label shown in UI
  value: string; // Actual tag value inserted into search (e.g. "rating:sensitive")
  category?: number | null; // 0=General, 1=Artist, 3=Copyright, 4=Character, 5=Meta
  post_count?: number;
  antecedent?: string | null; // Alias indicator (e.g. "miku" -> "hatsune_miku")
}

// types.ts

export interface DanbooruMediaVariant {
  type: "180x180" | "360x360" | "720x720" | "sample" | "original" | string;
  url: string;
  width: number;
  height: number;
  file_ext: string;
}

export interface DanbooruMediaAsset {
  id: number;
  md5: string;
  file_ext: string;
  file_size: number;
  image_width: number;
  image_height: number;
  status: string;
  file_key: string;
  is_public: boolean;
  variants?: DanbooruMediaVariant[];
}

export interface DanbooruPost {
  id: number;
  created_at: string;
  uploader_id: number;
  score: number;
  rating: "g" | "s" | "q" | "e";
  image_width: number;
  image_height: number;
  tag_string: string;
  tag_string_general: string;
  tag_string_character: string;
  tag_string_copyright: string;
  tag_string_artist: string;
  tag_string_meta: string;
  file_ext: string;
  file_size: number;

  // Restriction Flags
  is_banned: boolean;
  is_deleted: boolean;
  is_pending: boolean;
  is_flagged: boolean;

  // URLs (Nullable for restricted/banned content)
  file_url?: string | null;
  large_file_url?: string | null;
  preview_file_url?: string | null;

  // Modern Media Asset Object
  media_asset?: DanbooruMediaAsset | null;
}

export interface TagInputProps {
  tags?: string[];
  apiBaseUrl?: string;
  onsearch?: () => void;
}
