<script>
    import "./style.css"
    import {profile} from "../../store";
</script>
<h2>Профиль</h2>
<div>На данный момент функционал изменения настроек профиля закрыт</div>
<div class="card" style="margin-top: 30px; display: inline-block">
    <span style="font-weight: bold; margin-right: 20px">Ваше имя</span>
    <input value={$profile.name} style="width: 600px">
</div>
<br>
<div class="card" style="margin-top: 10px; display: inline-block">
    <span style="font-weight: bold; margin-right: 20px">Email&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>
    <input value={$profile.email} style="width: 600px">
</div>
<div></div>
<div class="card" style="margin-top: 10px; display: inline-block">
    <div style="font-weight: bold; display: inline-block; vertical-align: middle; margin-right: 20px">MqttToken/<br>HTTP API</div>
    <input value={$profile.mqttToken} style="width: 600px">
</div>