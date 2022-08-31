<script>
    import "./TopPanel.scss"
    import * as lrFetch from '../../common/lrFetch'
    import {profile} from "../../store"

    const onLogout = () => {
        lrFetch.get("/api/logout").then(response => {
            if (response.errno === 0) {
                profile.set({})
            }
        })
    }
</script>

<div class="top-panel">
    <span style="margin-right: 10px">{$profile.name}</span><img src="avatar.png" class="top-avatar">
    <div class="logout" title="Выйти" on:click={()=>onLogout()}><i class="fas fa-sign-out-alt"></i></div>
</div>
