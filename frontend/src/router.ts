import { createRouter } from "sv-router";

import type { SvelteComponent } from "svelte";
import {
  BookImage,
  GitPullRequest,
  Images,
  LayoutList,
  Palette,
  type LucideIcon,
} from "@lucide/svelte";

import MainPanel from "./lib/MainPanel.svelte";
import MediaPanel from "./lib/MediaPanel.svelte";
import FormPanel from "./lib/FormPanel.svelte";
import GitPushPanel from "./lib/GitPushPanel.svelte";
import KritaWorkspacePanel from "./lib/KritaWorkspacePanel.svelte";
import KritaWorkspaceItemPanel from "./lib/KritaWorkspaceItemPanel.svelte";
import KritaBoardPanel from "./lib/KritaBoardPanel.svelte";
import DanbooruPanel from "./lib/DanbooruPanel.svelte";
import { getDanbooruQuery } from "./stores/queryCacheStore";

export type RouteItem = {
  path: string;
  title: string;
  icon: LucideIcon;
  component: SvelteComponent | any;
  isVisible: boolean;
  closeSidebar?: boolean;
  keepQueryParam?: boolean;
  getQuery?: () => string;
};

function createRouterName(router: RouteItem[]) {
  const routerSvRouter = router.reduce(
    (acc, item) => {
      acc[item.path] = item.component;
      return acc;
    },
    {} as Record<string, any>,
  );

  return {
    routerSvRouter,
    routerName: router,
  };
}

const routerNameBuilder = createRouterName([
  {
    path: "/",
    title: "Dashboard",
    icon: LayoutList,
    component: MainPanel,
    isVisible: true,
  },
  {
    path: "/krita-workspace",
    title: "Krita Workspace",
    icon: Palette,
    component: KritaWorkspacePanel,
    isVisible: true,
  },
  {
    path: "/krita-workspace/:workspace",
    title: "Krita Workspace",
    icon: Palette,
    component: KritaWorkspaceItemPanel,
    isVisible: false,
  },
  {
    path: "/krita-workspace/:workspace/board/:board",
    title: "Krita Artwork",
    icon: Palette,
    component: KritaBoardPanel,
    isVisible: false,
    closeSidebar: true,
  },
  {
    path: "/danbooru",
    title: "Danbooru Reference",
    icon: BookImage,
    component: DanbooruPanel,
    isVisible: true,
    keepQueryParam: true,
    getQuery: () => {
      const res = getDanbooruQuery();
      return res;
    },
  },
  {
    path: "/media",
    title: "Media",
    icon: Images,
    component: MediaPanel,
    isVisible: true,
  },
  {
    path: "/create",
    title: "Create Post",
    icon: LayoutList,
    component: FormPanel,
    isVisible: false,
  },
  {
    path: "/git-push",
    title: "Git Push Artworks",
    icon: GitPullRequest,
    component: GitPushPanel,
    isVisible: true,
  },
]);

export const { routerName } = routerNameBuilder;
export const { p, route, navigate, isActive } = createRouter(
  routerNameBuilder.routerSvRouter,
);

export function getActiveRouteItem(): RouteItem | undefined {
  return routerName.find((item) => {
    const regexPattern = item.path
      .replace(/[.*+?^${}()|[\]\\]/g, "\\$&") // escape special regex chars
      .replace(/:[a-zA-Z0-9_]+/g, "[^/]+"); // replace :param with path segment regex

    return new RegExp(`^${regexPattern}$`).test(String(route.pathname));
  });
}

export function isSidebarClosed(): boolean {
  return Boolean(getActiveRouteItem()?.closeSidebar);
}
