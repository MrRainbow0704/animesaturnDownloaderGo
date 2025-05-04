import { mount } from "svelte"
import App from "$src/App.svelte"
import "$src/style.css"

export default mount(App, { target: document.getElementById("app")! });
