import { createRouter } from "sv-router";

import type { SvelteComponent } from "svelte";
import {
  GitPullRequest,
  Images,
  LayoutList,
  type LucideIcon,
} from "@lucide/svelte";

import MainPanel from "./lib/MainPanel.svelte";
import MediaPanel from "./lib/MediaPanel.svelte";
import FormPanel from "./lib/FormPanel.svelte";
import GitPushPanel from "./lib/GitPushPanel.svelte";

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
