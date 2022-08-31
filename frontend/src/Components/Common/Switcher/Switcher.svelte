<script>
    import "./Switcher.scss"
    export let value;
    export let initValue;
    export let action;
    $: initValue = parseInt(initValue);
    export let setValue;
    let clickSwitcher = () => {
        initValue = initValue ? 0 : 1;
        action(initValue)
    }
</script>

<div style="display: flex; justify-content: center; align-items: center;" >
    <div on:click={()=>clickSwitcher()} style="margin: 0 10px" class="switch-button {initValue ? ' active': ''}"><div class="switcher-label">{initValue ? "On" : "Off"}</div></div>
</div>