<script>
    import {onMount} from "svelte";
    import * as lrFetch from "../../common/lrFetch"
    import "./MainPage.scss";
    import DeviceCard from "../DeviceCard/DeviceCard.svelte";
    import DashCard from "../DashCard/DashCard.svelte";
    import Switcher from "../Common/Switcher/Switcher.svelte";

    let dashboard;

</script>
<h1>Dashboard</h1>
<button>+ Новый элемент</button>
<div style="margin: 10px 0">
    <DashCard color = "#fff" name="Карточка 1"><span style="color: #222; font-size: 34px">123</span></DashCard>
    <DashCard name="Кнопка"><button>Включить свет</button></DashCard>
</div>