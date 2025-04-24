import { useAuth } from "@/composeables/useAuth";
import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: () => import("../layouts/MainLayout.vue"),
      children: [
        {
          path: "",
          name: "home",
          component: HomeView,
          meta: {
            icon: "home",
            title: "Home",
          },
        },
        {
          path: "/about",
          name: "about",
          // route level code-splitting
          // this generates a separate chunk (About.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () => import("@/views/AboutView.vue"),
          meta: {
            icon: "info",
            title: "About",
          },
        },

        {
          path: "/mapextract",
          name: "mapextract",
          component: () => import("@/views/MapExtractView.vue"),
          meta: {
            icon: "map",
            title: "Map Extract",
          },
        },
        {
          path: "/localmapslist",
          name: "localmapslist",
          component: () => import("@/views/LocalMapsListView.vue"),
          meta: {
            icon: "map",
            title: "Local Maps List",
          },
        },
        {
          path: "/localmaps",
          name: "localmaps",
          component: () => import("@/views/PMMapsViewer.vue"),
          meta: {
            icon: "map",
            title: "Local Maps Viewer",
          },
        },
        {
          path: "/profile",
          name: "profile",
          component: () => import("@/views/ProfileView.vue"),
          meta: {
            icon: "person",
            title: "Profile",
            hideInNav: true,
          },
        },
      ],
    },
    {
      path: "/",
      component: () => import("../layouts/AnonymousLayout.vue"),
      children: [
        {
          path: "/login",
          name: "login",
          // route level code-splitting
          // this generates a separate chunk (About.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () => import("@/views/LoginView.vue"),
        },
        {
          path: "/register",
          name: "register",
          component: () => import("@/views/RegisterView.vue"),
        },
      ],
    },
  ],
});

//RouteGuard to check if user is authenticated and redirect to login page if not

// const unauthorizedRoutes = ["login", "register"];

// router.beforeEach(async (to, from, next) => {
//   const isAuthenticated = await useAuth().getAuthState();
//   if (!isAuthenticated && !unauthorizedRoutes.includes(to.name as string)) {
//     next({ name: "login" });
//   } else next();
// });

export default router;
