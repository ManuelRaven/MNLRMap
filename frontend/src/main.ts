// FILE: main.js

import { createApp } from "vue";
import { Quasar, Notify } from "quasar";
import quasarLang from "quasar/lang/de-DE";
import router from "./router";
// Import icon libraries
import "@quasar/extras/material-icons/material-icons.css";

// Import Quasar css
import "quasar/src/css/index.sass";

// Import your custom css
import "./assets/main.css";

// Assumes your root component is App.vue
// and placed in same folder as main.js
import App from "./App.vue";

const myApp = createApp(App);
myApp.use(router);
myApp.use(Quasar, {
  plugins: {
    Notify,
  }, // import Quasar plugins and add here
  lang: quasarLang,
});

// Assumes you have a <div id="app"></div> in your index.html
myApp.mount("#app");

// globals

myApp.config.globalProperties.GlobalAppName = "MNLRMap";

// Globals type
export {};

declare module "vue" {
  interface ComponentCustomProperties {
    GlobalAppName: string;
  }
}
