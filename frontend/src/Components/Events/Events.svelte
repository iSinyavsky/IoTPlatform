<style>
    .value{
        width: 100px;
    }
    .if-row{
        height: 80px;
        text-align: left;
        display: inline-flex;
        justify-content: space-between;
        align-items: center;
        position: relative;
    }

    .comparing{
        font-size: 30px;
        font-weight: bold;
        padding: 10px;
        cursor: pointer;
        color: #444;
    }

    .value{
        width: auto !important;
        display: inline-block;
        font-size: 30px;
        color: #4699d6;
    }
    .add-event{
        position: relative;
    }
    .content-card{
        position: relative;
    }
    .del {
        cursor: pointer;
        color: #666;
        position: absolute;
        top: 5px;
        right: 8px;
    }
    .trigger-style{
        color: #fff;
        font-size: 16px;
        background: #4699d6;
        padding: 20px;
        border-radius: 10px;

    }
    .disable{
        background: #2f324d !important;
    }
</style>

<script>
    import {openModal} from "../Modal/utils";
    import * as lrFetch from "../../common/lrFetch"
    import EventModal from "../Modal/EventModal/EventModal.svelte"
    import {variables} from "../../store";
    import {onMount} from "svelte";
    import * as common from "../../common";


    let if_event = {};
    let if_trigger = null;
    // if_event = {1: {value: "23", type: 1, operator: "="},
    //     7: {value: "23", type: 1, operator: "&ge;"}}
    let then_event = {};

    let addIf = (newIf) => {
        if_event = {...if_event, ...newIf}
        console.log(if_event)
    }
    let addTrigger = (newIf) => {
        if_event = {}
        if_trigger = newIf
        console.log(if_trigger)
    }
    let addThen = (newThen) => {
        then_event = {...then_event, ...newThen}
        console.log(then_event);
    }
    let events = [];
    onMount(()=>{
        getEvents();
    })

    function getEvents(){
        lrFetch.get("/api/getEvents").then(response => {
            events = response;
            console.log(events);
        })
    }

    function getVarById (id){
        return $variables.find((v)=>v.Id == id)
    }

    let saveTrigger = () => {
        lrFetch.post("/api/events", {if_event: if_event, then_event: then_event, if_trigger: if_trigger}).then(()=>{
            getEvents();
        })

       // location.reload();
    }
    let delEvent = (id) => {
        lrFetch.get("/api/removeEvent?id="+id)
        getEvents();
    }

    function getIntervalString (interval) {
        if (interval <= 60) {
            return "???????????? "+interval+" ??????????";
        } if (interval/60 < 24) {
            return "???????????? "+interval/60+" ????????";
        } else {
            return "???????????? "+interval/1440+" ????????";
        }
    }

    function timeLeft (time) {
        let now = new Date();
        time = new Date(time);
        let diffTime = time.getTime() - now.getTime();
        let s = 0;
        let m = 0;
        let h = 0;
        let d = 0;

        return secondsToHms(diffTime / 1000)
    }

    function secondsToHms(d) {
        d = Number(d);
        var h = Math.floor(d / 3600);
        var m = Math.floor(d % 3600 / 60);
        var s = Math.floor(d % 3600 % 60);

        var hDisplay = h > 0 ? h + (h == 1 ? " ??????, " : " ??????????, ") : "";
        var mDisplay = m > 0 ? m + (m == 1 ? " ????????????, " : " ??????????, ") : "";
        var sDisplay = s > 0 ? s + (s == 1 ? " ??????????????" : " ????????????") : "";
        return hDisplay + mDisplay + sDisplay;
    }
</script>
<h1>???????????????? (??????????????)</h1>
<div style="margin-top: -14px; margin-bottom: 20px">?????????????? ?????????????????? ???? ?????????? ???????????????? ?? ?????????????????? ?????????????? ??????????????????</div>

{#if $variables && Object.keys($variables).length !== 0}
<div class="content-card add-event">
    <div style="font-size: 20px; font-weight: bold; margin-bottom: 10px;">?????????? ??????????????</div>
    <div style="display: flex; align-items: center">
        <div style="margin-right: 20px; font-size: 40px; font-weight: bold">if</div>
        {#each Object.keys(if_event) as if_key, i}
            {#if i === 1}
                <div style="font-size: 14px; font-weight: bold; margin: 0 15px ">AND</div>
            {/if}
            <div class="if-row">
                <div>
                    <div class="device-card" style="min-width: 150px; min-height: 18px; font-size: 12px; padding: 15px; margin: 0; background: #1c1e33 !important;">{getVarById(if_key).Name}</div>
                </div>
                <div class="comparing" style="color: #333">{@html if_event[if_key].operator}</div>

                {#if if_event[if_key].type === 2}
                    <div>
                        <div class="device-card" style="min-width: 150px; min-height: 18px; font-size: 12px; padding: 15px; margin: 0; background: #1c1e33 !important;">{getVarById(if_event[if_key].value).Name}</div>
                    </div>
                    {:else}
                    <div class="value">
                        {if_event[if_key].value}
                    </div>
                {/if}
            </div>
        {/each}
        {#if if_trigger != null}
            <div class="if-row">
                <div class="trigger-style">
                    <i class="trigger-icon fas fa-clock"></i>
                {#if if_trigger.interval}
                    {getIntervalString(if_trigger.interval)}
                {:else if if_trigger.time}
                    ?????????? = {common.getNormalDateTime(if_trigger.time)}
                {/if}
                </div>
            </div>
        {/if}
        <button on:click={()=> openModal(EventModal, {addIf: addIf})} style="margin: 0 14px; background: limegreen">+ ?????? ????????</button>
        {#if if_trigger == null}
        <button on:click={()=> openModal(EventModal, {addTrigger: addTrigger, mode: 3})} style="margin: 0 14px; background: limegreen">?????? ???????????????? ?????????? <i class="trigger-icon fas fa-clock"></i>
        </button>
        {/if}
        <div></div>
    </div>
    <div style="display: flex; align-items: center">
        <div style="margin-right: 20px; margin-left:20px; font-size: 40px; font-weight: bold">then</div>
        {#each Object.keys(then_event) as then_key, i}
            {#if i === 1}
                <div style="font-size: 14px; font-weight: bold; margin: 0 15px ">AND</div>
            {/if}
            <div class="if-row">
                <div>
                    <div class="device-card" style="min-width: 150px; min-height: 18px; font-size: 12px; padding: 15px; margin: 0; background: #1c1e33 !important;">{getVarById(then_key).Name}</div>
                </div>
                <div class="comparing" style="color: #333">&#8594;</div>
                <div class="value">{then_event[then_key].value}</div>
            </div>
        {/each}
    <button on:click={()=> openModal(EventModal, {addThen: addThen, mode: 2})} style="margin: 0 14px; background: limegreen">+ ?????? ??????????</button>

    </div>

    <button style="margin: 20px 0 0 0" on:click={()=>saveTrigger()}>????????????????</button>
</div>


<h2>???????????? ??????????????????</h2>
    {#each events as event}
<div class={"content-card event-card " + (event.IsActive ? " disable" : "")} style="vertical-align: top; padding: 5px !important; width: 25%; color: #fff; background: #4699d6">
    <div class="del" style="color: #fff" on:click={()=>delEvent(event.Id)}>&#10006;</div>
    <div style="margin: 20px">
        {#if event.If_trigger.Time !== ""}
            <i class="trigger-icon fas fa-clock"></i>
        {/if}
        {#if event.If_trigger.Time !== "" && event.If_trigger.Interval === 0}
            ?? {common.getNormalDateTime(new Date(event.If_trigger.Time))}
        {/if}
        {#if event.If_trigger.Time !== "" && event.If_trigger.Interval !== 0}
            {getIntervalString(event.If_trigger.Interval)}
        {/if}
        {#if Object.keys(event.If_event).length > 0}<strong>????????</strong> ???????????????? ??????????????????{/if}
        {#each Object.keys(event.If_event) as ifEventKey}
        {#if getVarById(ifEventKey)}
            <strong>{getVarById(ifEventKey).Name}</strong>
            {:else }
            <span style="color: red">???????????????????? ??????????????</span>
        {/if}
           {@html event.If_event[ifEventKey].operator} {event.If_event[ifEventKey].value},
        {/each}
        {#if Object.keys(event.If_event).length > 0} <strong>??????????</strong> {/if}
        {#each Object.keys(event.Then_event) as thenEventKey}
            ?????????????????? ????????????????????

            {#if getVarById(thenEventKey)}
            <strong>{getVarById(thenEventKey).Name}</strong>
                {:else}
                <span style="color: red">???????????????????? ??????????????</span>
                {/if}
        ????????????????  {event.Then_event[thenEventKey].value}
            {/each}

    </div>
    {#if !event.IsActive}
    {#if event.If_trigger.Time !== ""}
        <div style="font-size: 14px; margin: -10px 0 0 20px; display: inline-block; max-width: 80%; background: #4659e0; color: #fff; padding: 5px">???????????????????? ?????????? {timeLeft(event.If_trigger.Time)}</div>
    {/if}
        {:else}
        <div style="font-size: 14px; margin: -10px 0 0 20px; display: inline-block; max-width: 80%; background: orangered; color: #fff; padding: 5px">??????????????????</div>

    {/if}
    <div style="font-size: 12px; margin: 5px; text-align: right">?????????????? {common.getNormalDateTime(event.CreationDate)}</div>

</div>
    {/each}



{/if}
