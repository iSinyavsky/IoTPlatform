<script>

	import Router from 'svelte-spa-router'
	import {link} from 'svelte-spa-router'
	import * as lrFetch from './common/lrFetch'
	import MainPage from "./Components/MainPage/MainPage.svelte";
	import Events from "./Components/Events/Events.svelte"
	import Variables from "./Components/Variables/Variables.svelte"
	import VariableInfo from "./Components/VariableInfo/VariableInfo.svelte"
	import Sidebar from "./Components/Sidebar/Sidebar.svelte";
	import Profile from "./Components/Profile/Profile.svelte"
	import {BASE_URL, FRONT_DOMAIN, PROTOCOL} from "./consts/consts";
	import {onMount} from "svelte";
	import LoginPage from "./Components/LoginPage/LoginPage.svelte";
	import {profile, variables, variableValues} from "./store"
	import TopPanel from "./Components/TopPanel/TopPanel.svelte"
	import CodeGen from "./Components/CodeGen";
	import * as ws from "./websocket"
	const routes = {
		// Exact path
		'/': MainPage,
		'/variables': Variables,
		'/variables/:id': VariableInfo,
		'/events': Events,
		"/codeGen": CodeGen,
		"/profile": Profile

	};

	let isLogin = false;
	let isAuth = false;

	onMount(() => {

		if (window.location.hash.split("access_token=").length > 1){
			let token = window.location.hash.split("access_token=")[1].split("&")[0]
			lrFetch.post("/api/yandex/saveYandexToken", {token}).then(response => {

			})
		}
		lrFetch.get("/api/getActor").then(response => {
			if (response.errno === 0) {
				profile.set(response)

				ws.init();
				variables.getVariables()


			} else {
				profile.set({})
			}
		})


	})



</script>
{#if $profile == null}
<div>Loading...</div>
{:else if Object.keys($profile).length >0}
<Sidebar/>
<main>
	<TopPanel></TopPanel>
	<Router {routes}/>
</main>
{:else}
	<LoginPage/>
{/if}