<script>
    import * as lrFetch from '../../common/lrFetch'
    import {profile} from "../../store"
    import "./LoginPage.scss"
    let email = "";
    let password = "";

    let onLogin = () => {
        lrFetch.post("/api/login", {email, password}).then(response=>{
            profile.set({userId: response.userId, mqttToken: response.mqttToken})
        });
    }

    let onRegister = () => {
        lrFetch.post("/api/register", {email, password}).then(response=>{
            profile.set({userId: response.userId, mqttToken: response.mqttToken})
        });
    }

    let loginOrRegister = 0; //0-login, 1 - register

</script>

{#if loginOrRegister === 0}
<div class="login-form">
    <div style="text-align: center">Авторизация</div>
    <div class="field input-field">
        <div class="label">E-mail</div>
        <i class="fas fa-user" style="margin-right: 10px"></i>
        <input type="text" bind:value={email} placeholder="Type email">
    </div>
    <div class="field input-field">
        <div class="label">Password</div>
        <i class="fas fa-key" style="margin-right: 10px"></i>
        <input type="password" bind:value={password} placeholder="Type password">
    </div>
    <div class="field">
        <button on:click={()=> onLogin()}>Войти</button>
    </div>
    <div class="field" style="text-align: center; font-size: 12px; cursor: pointer">Нет аккаунта? <span style="color: #4699d6" on:click={()=>loginOrRegister = 1}>Регистрация</span></div>
</div>

{:else}
    <div class="login-form">
        <div style="text-align: center">Регистрация</div>
        <div class="field input-field">
            <div class="label">E-mail</div>
            <i class="fas fa-user" style="margin-right: 10px"></i>
            <input type="text" bind:value={email} placeholder="Type email">
        </div>
        <div class="field input-field">
            <div class="label">Password</div>
            <i class="fas fa-key" style="margin-right: 10px"></i>
            <input type="password" bind:value={password} placeholder="Type password">
        </div>
        <div class="field">
            <button on:click={()=>onRegister()}>Регистрация</button>
        </div>
        <div class="field" style="text-align: center; font-size: 12px; cursor: pointer">Есть аккаунт? <span style="color: #4699d6" on:click={()=>loginOrRegister = 0}>Авторизация</span></div>
    </div>
{/if}