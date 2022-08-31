<script>
    import "./VariableInfo.scss"
    import LineChart from "../Charts/LineChartWithZoom/LineChartWithZoom.svelte"
    import LineChartWithZoom from "../Charts/LineChartWithZoom/LineChartWithZoom.svelte";
    import Switcher from "../Common/Switcher/Switcher.svelte"
    import RangeSlider from "../Common/RangeSlider/RangeSlider.svelte"
    import {openModal} from "../Modal/utils";
    import ConnectGuide from "../Modal/ConnectGuide/ConnectGuide.svelte"
    export let params = {}
    import {variables, profile, variableValues} from "../../store";
    import * as lrFetch from '../../common/lrFetch'
    import {onMount} from "svelte";
    import HTTPGuide from "../Modal/HTTPGuide/HTTPGuide.svelte";

    let generateData = () => {
        var chartData = [];
        var firstDate = new Date();
        firstDate.setDate(firstDate.getDate() - 1000);
        var value = 100;
        for (var i = 0; i < 50; i++) {
            // we create date objects here. In your data, you can have date strings
            // and then set format of your dates using chart.dataDateFormat property,
            // however when possible, use date objects, as this will speed up chart rendering.
            var newDate = new Date(firstDate);
            newDate.setDate(newDate.getDate() + i);

            value += Math.round((Math.random()<0.5?1:-1)*Math.random()*10);

            chartData.push({
                date: newDate,
                value: value
            });
        }
        return chartData;
    }

    let deviceActions;
    let deviceServiceInfo;
    onMount(()=>{
        variables.subscribe(vars => {
            for(let i=0; i<vars.length; i++){
                if(vars[i].Id === parseInt(params.id)) {
                    variable = {...vars[i]};

                    if(variable.ServiceName !== "") {
                        lrFetch.get("/api/yandex/getDevicesActions").then(response => {
                            deviceActions = response;
                            console.log(deviceActions);
                        })
                        lrFetch.get("/api/yandex/getDevice?deviceID="+variable.IntID).then(response => {
                            deviceServiceInfo = response;
                            console.log(deviceServiceInfo);
                        })
                    }

                    return
                }
            }
        })

        variableValues.getVariableValues(params.id);
    })

    let variable = null;
    let mqttToken = "";
    let newValue;
    profile.subscribe(value=>{
        mqttToken = value.mqttToken;
    })

    variableValues.subscribe(varValues => {
        if(params.id)
            if(varValues[params.id]){
                if(varValues[params.id].length > 0)
                    newValue = varValues[params.id][varValues[params.id].length-1].Current_Value
                else newValue = "Пусто";
            }
    })


    const save = () => {
        lrFetch.put("/api/updateVariable?id="+params.id, {name: variable.Name, label: variable.Label})
    }

    const setValue = (value) => {
        lrFetch.post("/api/setValue", {varId: parseInt(params.id), value: value, mqttToken: mqttToken, label: variable.Label})
    }

    const onOffYandex = (variableId, varID, type,value) => {
        lrFetch.get("/api/yandex/sendValue?varID="+variableId+"&deviceID="+varID+"&type="+variable.Capability+"&value="+value+"&deviceType="+variable.IntType).then(()=>{
        })
    }

    const rangeYandex = (variableId, varID, type,value) => {
        lrFetch.get("/api/yandex/sendValue?varID="+variableId+"&deviceID="+varID+"&type="+variable.Capability+"&value="+value+"&deviceType="+variable.IntType).then(()=>{
        })
    }

    console.log("aa", localStorage.getItem("showChart"+params.id));

    let isShowChart = localStorage.getItem("showChart"+params.id) !== null  ? parseInt(localStorage.getItem("showChart-"+params.id)) : true;
    const showChart = (status) => {
        isShowChart = status
        localStorage.setItem("showChart"+params.id, ""+status);
        console.log("ee", localStorage.getItem("showChart"+params.id));
    }

    const exportCsv = () => {
        let csvContent = "csv,csv"
        var encodedUri = encodeURI(csvContent);
        var link = document.createElement("a");
        link.setAttribute("href", encodedUri);
        link.setAttribute("download", "my_data.csv");
        document.body.appendChild(link); // Required for FF

        link.click();
    }

</script>

{#if variable}
<h2>{variable.Name}</h2>
<div style="margin: 30px 0" class="row">
    <div class="col-2">
        <div class="card">
            <div  class="card-label">Название</div>
            <input on:blur={()=>save()} bind:value={variable.Name} class="title-input main">
<!--            <div  class="card-label">Описание</div>-->
<!--            <textarea class="title-input" style="resize: none;"></textarea>-->
            {#if variable.ServiceName === ""}
                <div  class="card-label">Label (for mqtt)</div>
                <input on:blur={()=>save()} class="title-input" bind:value={variable.Label}/>
                {:else}
                <div class="card-label">Integration DeviceID</div>
                <input readonly class="title-input" bind:value={variable.IntID}/>
                <div class="card-label">Service</div>
                <input readonly class="title-input" bind:value={variable.ServiceName}/>
            {/if}

        </div>
        {#if variable.ServiceName !== ""}
            {#if variable.Capability === "devices.capabilities.on_off"}
            <div class="card" style="margin-top: 10px; width: 260px; display: inline-block">
                <div class="card-label">Переключатель</div>
                <Switcher initValue={newValue} action={(value)=>onOffYandex(variable.Id, variable.IntID, variable.Capability,value)}></Switcher>
            </div>
            {/if}
            {#if variable.Capability === "devices.capabilities.range"}
                <div class="card" style="margin-top: 10px; width: 260px; display: inline-block">
                    <div class="card-label">Слайдер</div>
                    <RangeSlider value={newValue} action={(value)=>rangeYandex(variable.Id, variable.IntID, variable.Capability,value)}></RangeSlider>
                </div>
            {/if}
        {/if}
        {#if variable.ServiceName === "" && isShowChart}
        <div class="card" style="margin-top: 10px">
            <div class="card-label">Управление (произвольное)</div>
            <div><input bind:value={newValue} class="title-input" ></div>
            <button style="margin: 10px; width: calc(100% - 20px);" on:click={()=>setValue(newValue)}>Сохранить</button>
        </div>
        <div class="card" style="margin-top: 10px">
            <div class="card-label">Управление (переключатель)</div>
            <Switcher initValue={newValue} action={(value)=>setValue(value+"")}></Switcher>
        </div>
        <div class="card" style="margin-top: 10px">
            <div class="card-label">Управление (слайдер)</div>
            <RangeSlider></RangeSlider>
            <button style="margin: 10px; width: calc(100% - 20px);">Сохранить</button>
        </div>
        <div style="margin: 10px auto; text-align: center; color: #4699d6; font-size: 14px; text-decoration: underline; cursor: pointer">+ Добавить виджет управления</div>
        {/if}
    </div>
    <div class="col-10" style="position: relative">
        {#if true}
        <div class="chart-panel">
            {#if isShowChart}
                <button style="background: red" on:click={()=>showChart(0)}>Не показывать график</button>
            {:else}
                <button style="background: red" on:click={()=>showChart(1)}>Показывать график</button>
            {/if}
            <button on:click={()=>openModal(HTTPGuide)}>HTTP гайд</button>

            <button on:click={()=>openModal(ConnectGuide)}>MQTT гайд</button><button on:click={()=>exportCsv()}>Экспорт в CSV</button><button>Экспорт в JSON</button>
        </div>
        {#if isShowChart}
            <div class="card">
                {#if $variableValues[params.id] && $variableValues[params.id].length > 0}
                    <LineChartWithZoom style="height: 600px; margin-top: -20px" data={$variableValues[params.id]}></LineChartWithZoom>
                {/if}
            </div>
            {:else}
                <div class="card col-3" style="margin-top: 10px">
                    <div class="card-label">Управление (произвольное)</div>
                    <div><input bind:value={newValue} class="title-input" ></div>
                    <button style="margin: 10px; width: calc(100% - 20px);" on:click={()=>setValue(newValue)}>Сохранить</button>
                </div>
                <div class="card col-3" style="margin-top: 10px">
                    <div class="card-label">Управление (переключатель)</div>
                    <Switcher initValue={newValue} setValue={setValue}></Switcher>
                </div>
                <div class="card col-3" style="margin-top: 10px">
                    <div class="card-label">Управление (слайдер)</div>
                    <RangeSlider></RangeSlider>
                    <button style="margin: 10px; width: calc(100% - 20px);">Сохранить</button>
                </div>
        {/if}
        {:else}
            {#if deviceServiceInfo}
            <div>Это внешнее устройство, статистика в реальном времени недоступна.</div>
            {#each deviceServiceInfo.capabilities as capabilitiy}
                {#if capabilitiy.type === "devices.capabilities.on_off"}
            <div class="card" style="margin-top: 10px; width: 300px; display: inline-block">
                <div class="card-label">Управление (переклюк)</div>
                <Switcher value={capabilitiy.state.value} action={(value)=>onOffYandex(variable.IntID, capabilitiy.type,value)}></Switcher>
            </div>
                {/if}
                {#if capabilitiy.type === "devices.capabilities.range"}
            <div class="card" style="margin-top: 10px; width: 300px; display: inline-block">
                <div class="card-label">Управление (слайдер)</div>
                <RangeSlider value={capabilitiy.state.value} action={(value)=>rangeYandex(variable.IntID, capabilitiy.type,value)}></RangeSlider>
            </div>
                {/if}
            {/each}
                {:else}
                <div>Загрузка</div>
            {/if}
        {/if}
    </div>
</div>
{/if}