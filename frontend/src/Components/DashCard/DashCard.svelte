<script>
    import "./style.scss"
    export let name;
    export let color;
    export let value;
    export let id;
    export let serviceName;
    export let isDeleted = false;
    import * as lrFetch from '../../common/lrFetch'

    const deleteVariable = (e, id) => {
        e.stopPropagation();
        lrFetch.get("/api/removeVariable?id="+id).then(()=>{
            isDeleted = true;
        })
    }

    const showVariable = (href) => {
        location.href = href;
    }
</script>

{#if !isDeleted}
    <div on:click={()=>showVariable("#/variables/"+id)} class="device-card" style={"color: #222; background: #fff"}>
        <div class="title-section"><div class="name"> {name}</div></div>
        <div style="margin: 20px; text-align: center"><slot></slot></div>
        <div class="delete" on:click={((e)=>deleteVariable(e, id))}>&#10006;</div>
    </div>
{/if}