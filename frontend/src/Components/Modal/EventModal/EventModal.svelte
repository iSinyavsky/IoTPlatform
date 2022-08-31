<script>
 import Modal from "../index.svelte"
 import DatePicker from "@beyonk/svelte-datepicker/src/components/DatePicker.svelte";
 import  "./whatislr.scss";
 import {variables} from "../../../store"
 import * as common from "../../../common";
 export let addIf;
 export let addThen;
 export let addTrigger;
 let visible = true;
 let choseIf = null;
 let ifOperator = null;
 let ifValue = "";
 let ifValueId = null;

 let thenValue = "";

 let thenVariable;
 let ifContainer;
 let showIfValueVariables = false;

 let operators = ["&lt;", "&le;", "=", "&ge;", "&gt;"];

 function getVarById (id){
     return $variables.find((v)=>v.Id === id)
 }

 function setIf (id) {
     choseIf = id
     setTimeout(()=>{
        ifContainer.scrollTo({top: 10000})
     },200)
 }

 function setThen(id) {
    thenVariable = id;
     setTimeout(()=>{
         ifContainer.scrollTo({top: 10000})
     },200)
 }

 function setOperator(o) {
     ifOperator=o;
     setTimeout(()=>{
         ifContainer.scrollTo({top: 10000})
     },200)

 }

 function scrollToEnd(){
     setTimeout(()=>{
         ifContainer.scrollTo({top: 10000})
     },200)
 }

 export let mode = 1;
 let if_event = {};
 let then_event = {};
 let if_trigger = {};

 let isInterval = false;
 let intervalScale = "60";
 let intervalValue = 2;

 function saveIfTrigger () {
     let intervalResult = 0;
     if (isInterval) {
        intervalResult = intervalValue * intervalScale
     }
     if_trigger = {
     }
     if (selected) {
         if_trigger["time"] = selected.date
     }
     if (intervalResult !== 0) {
         if_trigger["interval"] = intervalResult
     }
     addTrigger(if_trigger)
     visible = false;
 }

 function saveIf(){


    if_event[choseIf] = {
        value: ifValueId !== null && ifValueId || ifValue,
        type: ifValueId !== null && 2 || 1,
        operator: ifOperator
    }
    addIf(if_event);
    visible = false;
 }
 let selected = undefined;

 function saveThen(){
     then_event[thenVariable] = {
         value: thenValue
     }
     addThen(then_event);
     visible = false;
 }


</script>

<style>
    .comparing-block{
        display: flex;
        text-align: center;
        width: 550px;
        margin: 60px auto;
        justify-content: space-between;
    }
    .comparing{
        font-size: 60px;
        font-weight: bold;
        padding: 10px;
        cursor: pointer;
        color: #1c1e33;
    }
    .comparing.active{
        color: #4699d6;
    }
    .comparing:hover{
        transform: scale(1.1);
        color: #4699d6;
    }
    .head{
        margin: 0 auto 0 auto;
        width: 500px;
        text-align: left;
        display: flex;
        justify-content: space-between;
        align-items: center;

    }
    .value{
        font-size: 30px;
        color: #4699d6;
    }
    .ifvalue{
        width: 140px;
        border: none;
        border-bottom: 2px solid #4699d6;
        font-size: 40px;
        outline: none;
    }

    .if-container{
        padding-top: 20px;
        border-top: 2px solid #4699d6;
    }

    .device-card {
        font-size: 12px;
        min-width: 140px;
        min-height: 40px;
        background: #1c1e33 !important;
    }
    .device-card.active{
        background: #4699d6 !important;
    }

    .interval-button{
        background: #4699d6;
        color: #fff;
        display: inline-block;
        padding: 10px 20px;
        margin: auto;
        border: 1px solid #eee;
        text-align: center;
        width: 250px;
        text-decoration: none;
        cursor: pointer;
        font-size: 12px;
        border-radius: 7px;
        box-shadow: 0px 0px 3px rgb(0 0 0 / 10%);
        line-height: 27px;
    }
    .interval-form{
        padding: 20px;
        width: 60%;
        margin: 10px auto;
        border-radius: 10px;
        background: #2f324d;
        color: #fff;
    }
    .external{
        position: relative;
        background: #2f324d !important;
    }
    .external:after{
        position: absolute;
        bottom: 5px;
        right: 5px;
        content: "yandex";
        font-size: 10px
    }

</style>
{#if visible}
<Modal width="800">

    {#if mode === 1}
    {#if choseIf !== null}
    <div class="head">

            <div style="marin-right: 30px; font-size: 40px; font-weight: bold">if</div>
            <div>
                <div style="margin: 0 10px; color: #444; font-size: 12px">Последнее значение</div>
                <div class="device-card" style="min-width: 180px; min-height: 20px">{getVarById(choseIf).Name}</div>
            </div>
            {#if ifOperator !== null}
                <div class="comparing" style="color: #333">{@html ifOperator}</div>
            {/if}
            {#if ifValueId}
                <div class="device-card" style="min-width: 180px; min-height: 20px">{getVarById(ifValueId).Name}</div>
                {:else}
                <div class="value">{ifValue}</div>
            {/if}
            <button style="margin-top: 3px" on:click={()=>saveIf()}>Готово</button>

    </div>
    {/if}
    <div bind:this={ifContainer} class="if-container" style="text-align: center; max-height: 500px; overflow-y: scroll">
        <div style="font-size: 20px; margin-bottom: 20px">Выберите переменную для сравнения</div>
        {#each $variables as variable}
            <div on:click={()=>setIf(variable.Id)} class={(choseIf === variable.Id && "device-card event-card active" || "device-card") + (variable.ServiceName !== '' ? " external" : "")}>{variable.Name}</div>
        {/each}
        {#if choseIf !== null}
            <div style="font-size: 20px; margin: 20px 0">Выберите оператор сравнения</div>
            <div class="comparing-block">
                {#each operators as o}
                    <div class={ifOperator === o && "comparing active" || "comparing"} on:click={()=>setOperator(o)}>{@html o}</div>
                {/each}
            </div>
        {/if}
        {#if choseIf !== null && ifOperator !== null}
            <div style="font-size: 20px; margin: 20px 0" class="if-value">С чем сравнивать</div>
            <div>Введите значение <input class="ifvalue" bind:value={ifValue}>,</div>
            <div>или выберите переменную для сравнения <button on:click={()=>{showIfValueVariables = true; scrollToEnd()}}>Показать переменные </button></div>

            {#if showIfValueVariables}
                {#each $variables as variable}
                    <div class="device-card" on:click={()=>ifValueId = variable.Id} style="">{variable.Name}</div>
                {/each}
            {/if}
        {/if}
    </div>
    {/if}

    {#if mode === 2}
        <div bind:this={ifContainer} class="if-container" style="text-align: center; max-height: 500px; border-top: none; overflow-y: scroll">
            <div style="font-size: 20px; margin-bottom: 20px">Выберите переменную, у которой поменяется значение</div>
            {#each $variables as variable}
                <div on:click={()=>setThen(variable.Id)} class={(choseIf === variable.Id && "device-card active" || "device-card")+ (variable.ServiceName !== '' ? " external" : "")} style="">{variable.Name}</div>
            {/each}

            {#if thenVariable}
            <div style="font-size: 20px; margin: 20px 0" class="if-value">Присвоить значение</div>
            <div>Введите значение <input class="ifvalue" bind:value={thenValue}></div>
            {/if}
        </div>

        {#if thenVariable}
            <div style="text-align: right"><button on:click={()=>saveThen()}>Готово</button></div>
        {/if}
    {/if}

    {#if mode === 3}
        <div style="font-size: 20px; margin-bottom: 20px">Время срабатывания триггера</div>
        <DatePicker time={true} continueText="Выбрать время"
                    placeholder={"По определенному времени"} format="DD-MM-YYYY H:mm" on:date-selected={(e) => selected = e.detail}/>
        <div class="interval-button" on:click={()=>{isInterval = true; selected = undefined}}>По интервалу (например каждый час)</div>
        {#if !selected && isInterval}
            <div class="interval-form">
                Каждые <input style="width: 100px" bind:value={intervalValue} placeholder="2">
                <select bind:value={intervalScale}>
                    <option value="1">минут</option>
                    <option value="60">часа</option>
                    <option value="1440">дня</option>
                </select>
            </div>
        {/if}
        {#if selected}
        <div style="margin: 20px auto; font-size: 20px; color: #fff; background: #2f324d; width: 60%; padding: 30px; box-sizing: border-box;">
            Триггер выполнится в {common.getNormalDateTime(selected.date)}
        </div>
        {/if}
        {#if selected || isInterval}
            <div style="text-align: right" on:click={()=>saveIfTrigger()}><button>Готово</button></div>
        {/if}
    {/if}
</Modal>
{/if}