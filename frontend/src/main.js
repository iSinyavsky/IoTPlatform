import App from './App.svelte';
import "./App.scss"
const app = new App({
	target: document.body,
	hydratable: true
});

window.app = app;

export default app;