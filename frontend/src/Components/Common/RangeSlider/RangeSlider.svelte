<script>
    import "./RangeSlider.scss"
    export let enable = false;
    export let action;
    export let value = 0;
    const changeSlider = () => {
        console.log(value);
        action(value);
    }
</script>

<div style="display: flex; justify-content: center; align-items: center; position:relative;" >
    <div class="slider-value">{value}</div>
    <input type="range" min="0" max="100" class="slider" bind:value={value} on:change={()=>changeSlider()}>
</div>