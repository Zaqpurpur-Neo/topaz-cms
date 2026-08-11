import { createRouter } from "sv-router";

import type { SvelteComponent } from "svelte";
import {
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

type RouteItem = {
  path: string;
  title: string;
  icon: LucideIcon;
  component: SvelteComponent | any;
  isVisible: boolean;
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
