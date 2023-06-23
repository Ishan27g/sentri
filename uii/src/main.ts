
import { createApp } from 'vue'
import PrimeVue from "primevue/config";
import App from './App.vue'

import Button from "primevue/button";
import Card from "primevue/card";
import ScrollPanel from "primevue/scrollpanel";
import SelectButton from "primevue/selectbutton";
import Textarea from "primevue/textarea";
import Dialog from "primevue/dialog";



import Listbox from 'primevue/listbox';
import VirtualScroller from 'primevue/virtualscroller';

import "primeflex/primeflex.css";
import "primevue/resources/themes/mdc-dark-indigo/theme.css";


const app = createApp(App)


app.use(PrimeVue);

app.component("Button", Button);
app.component("Card", Card);
app.component("ScrollPanel", ScrollPanel);
app.component("SelectButton", SelectButton);
app.component("Textarea", Textarea);
app.component("Dialog", Dialog);


app.component("Listbox", Listbox);
app.component("VirtualScroller", VirtualScroller);


app.mount('#app')
