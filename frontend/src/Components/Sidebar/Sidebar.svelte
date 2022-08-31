<script>
 import "./style.scss";

 let loginForm = false;
 let registerForm = false;
 let userForm = false;
 import { fly } from 'svelte/transition';

 export let hide = false;
 let miniSidebar = true



</script>
<style>
</style>
{#if !hide}
<div class={miniSidebar ? "sidebar hide" : "sidebar"} in:fly="{{x: -200, duration: 300}}" out:fly="{{x: -200, duration: 300}}">
    <div class="sidebar-head">
        <div style="text-align: center; font-weight: bold; font-size: 10px">mThings</div>
    </div>
    <div class="menu">
        <div class="menu-item"><a href="#/"><i class="fas fa-tachometer-alt"></i><div class="item-title">Главная</div> </a></div>
        <div class="menu-item"><a href="#/variables"><i class="fas fa-microchip"></i><div class="item-title">Устройства</div></a></div>
        <div class="menu-item"><a href="#/events"><i class="fas fa-tools"></i><div class="item-title">Условия</div></a></div>
        <div class="menu-item"><a href="#/codeGen"><i class="fas fa-code"></i><div class="item-title">CodeGen</div></a></div>
        <div class="menu-item"><a href="#/profile"><i class="fas fa-user"></i><div class="item-title">Профиль</div></a></div>
    </div>
    <div class="hide-button" on:click={()=>miniSidebar = !miniSidebar}>&larr;&rarr;</div>
</div>
{/if}
