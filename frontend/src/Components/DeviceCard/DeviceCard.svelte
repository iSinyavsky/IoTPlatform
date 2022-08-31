<script>
    import "./style.scss"

    export let name;
    export let color;
    export let value;
    export let style;
    export let id;
    export let serviceName;
    export let isDeleted = false;
    import * as lrFetch from '../../common/lrFetch'
    import SettingsComponent from "./SettingsComponent/SettingsComponent.svelte";
    import {afterUpdate} from "svelte";

    const deleteVariable = (e, id) => {
        e.stopPropagation();
        lrFetch.get("/api/removeVariable?id=" + id).then(() => {
            isDeleted = true;
        })
    }

    const showVariable = (href) => {
        location.href = href;
    }

    let settingsComponent = false;

    if (!style) {
        style = {
            bg: "#4699d6",
            icon :"fa-microchip"
        }
    } else {

    }
    afterUpdate(()=>{
        if (!style) {
            style = {
                bg: "#4699d6",
                icon :"fa-microchip"
            }
        } else {

        }
    })

    const hideSettingsComponent = () =>{
        settingsComponent = false;
    }

</script>

{#if !isDeleted}
    <div on:click={()=>showVariable("#/variables/"+id)} class="device-card" style={`background: ${style ? style.bg : "#4699d6"};`}>
        {#if settingsComponent}<SettingsComponent hideSettingsComponent={hideSettingsComponent} id={id} style={style}></SettingsComponent>{/if}
        <div class="title-section"><i class={"fas "+(style ? style.icon : "fa-microchip")} style="margin-right: 10px"></i><div class="name"> {name}</div></div>
        <div class="value">{value ? value : "Пусто"}</div>
        <div class="service">{#if serviceName}{serviceName}{/if}</div>
        <div class="var-settings" on:click={(e)=>{e.stopPropagation(); settingsComponent = !settingsComponent}}>&#9881;</div>
        <div class="delete" on:click={((e)=>{ deleteVariable(e, id)})}>&#10006;</div>
    </div>
{/if}