<style>
    .v-settings{
        background: #fff;
        position: absolute;
        right: -30px;
        top: 30px;
        min-height: 40px;
        width: 180px;
        box-shadow: 0 4px 4px rgba(0,0,0,0.1);
        padding: 10px;
        border-radius: 6px;
        color: #222;
        z-index: 10000;
    }
    .title{
        font-size: 15px;
        text-align: center;
        margin-bottom: 10px;
    }
    .color:hover{
        opacity: .7;
    }
    .color{
        width: 40px;
        height: 40px;
        margin: 2px;
        display: inline-block;
        border-radius: 6px;
    }
    .icon:hover{
        opacity: .7;
    }
    .icon{
        display: inline-block;
        width: 40px;
        height: 40px;
        margin: 2px;
        font-size: 24px;
        text-align: center;
    }

</style>
<script>
    import { clickOutside } from '../../../common/index';
    import * as lrFetch from "../../../common/lrFetch"
    import {params} from "../../VariableInfo/VariableInfo.svelte";
    import {variables} from "../../../store";
    let visible = true;
    let colors = ["#4699d6", "#32C12C", "#D40C00", "#50342C", "#444444", "#84C144", "#798CD2", "#131F39"];
    let icons = ["fa-microchip", "fa-lightbulb", "fa-thermometer-empty", "fa-toggle-on", "fa-fan", "fa-pump-soap",  "fa-sun", "fa-bell"];
    export let style;
    export let id;
    export let hideSettingsComponent;

    const save = (e,param, value) => {
        e.stopPropagation();
        style[param] = value;
        lrFetch.post("/api/updateStyleVariable?id="+id, style).then(()=>{
            variables.getVariables()
        })
    }
</script>

{#if visible}
<div class="v-settings"
     use:clickOutside
     on:click_outside={hideSettingsComponent}
>
    <div class="title">Стилизация карточки</div>
    {#each colors as color}
        <div class="color" style={"background: "+color} on:click={(e)=>save(e, "bg", color)}></div>
    {/each}
    <div style="margin: 10px 0;">
    {#each icons as icon}
        <div class="icon"  on:click={(e)=>save(e,"icon", icon)}><i class={"fas "+icon} style="margin-right: 10px"></i></div>
    {/each}
    </div>
</div>
{/if}