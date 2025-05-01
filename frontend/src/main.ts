import {mount} from "svelte"
import App from "./App.svelte"
import "./style.css"

export default mount(App, { target: document.getElementById("app")! });
