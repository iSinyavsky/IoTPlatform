<script>
    import DeviceCard from "../DeviceCard/DeviceCard.svelte"
    import {onMount} from "svelte";
    import * as lrFetch from '../../common/lrFetch'
    import {variables, profile} from "../../store"
    import * as common from "../../common";
    import {FRONT_DOMAIN} from "../../consts/consts";
    const addVariable = () => {
        let label = common.MD5(""+new Date().getTime());
        lrFetch.post("/api/addVariable", {"name": "Ещё одна переменная", "label":  label}).then(response => {
            variables.getVariables()
            location.reload();
        })

    }

    let lastValues = [];
    onMount(()=>{
        lrFetch.get("/api/getLastValues").then(response=>{
            lastValues = response;
        })
    })

    const findValueByVarId = (array, varId) => {
        return array.find((el)=>el.Id === varId);
    }

    let mergeYandex = () => {
        if ($profile.yaToken === "") {
            window.open("https://oauth.yandex.ru/authorize?response_type=token&client_id=070022eeb1dc4403b1fd97507ec2884d&redirect_uri=https://"+FRONT_DOMAIN+"/#/yandex_login").focus();
        } else {
            lrFetch.get("/api/yandex/getDevices").then((response)=>{
                console.log(response)
            })
        }

    }

    const getVariablesByServiceFilter = (variables, isService = false) => {
        if (isService){
            return variables.filter((el)=>el.ServiceName !== "")
        }
        return variables.filter((el)=>el.ServiceName === "")
    }
</script>

<h1>Датчики и исполнительные устройства</h1>
<button on:click={()=>addVariable()}>+ Добавить переменную</button>
<button class="white-button" style="margin: 0 10px" on:click={()=>mergeYandex()}>Объеденить с Yandex Home</button>
<h2>Ваши устройства</h2>
<div style="margin: 10px 0">
    {#if $variables.length > 0}
        {#each getVariablesByServiceFilter($variables) as variable}
            <DeviceCard id={variable.Id} name={variable.Name} value={lastValues.length != 0 && findValueByVarId(lastValues, variable.Id).Value} style={variable.Style} color="#4699d6"></DeviceCard>
        {/each}
    {/if}
</div>
<h2>Сторонние устройства</h2>
<div style="margin: 10px 0">
    {#if $variables.length > 0}
        {#each getVariablesByServiceFilter($variables, true) as variable}
            <DeviceCard id={variable.Id} name={variable.Name} serviceName={variable.ServiceName} value={lastValues.length != 0 && findValueByVarId(lastValues, variable.Id).Value} style={variable.Style} color="#1c1e33"></DeviceCard>
        {/each}
    {/if}
</div>