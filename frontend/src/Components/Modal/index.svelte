<svelte:window on:keyup={handleKeyup}/>

<script>
    import './modal.sass';
    import { createEventDispatcher } from 'svelte';

    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';

    import { clickOutside } from '../../common/index';

    export let width = 280;
    export let additionalClass = '';
    export let visible = false;

    const dispatch = createEventDispatcher();
    const onBeforeShowModal = () => dispatch('onBeforeShowModal');
    const onAfterShowModal = () => {dispatch('onAfterShowModal')};
    const onBeforeCloseModal = () =>  dispatch('onBeforeCloseModal');
    const onAfterCloseModal = () => dispatch('onAfterCloseModal');
    const onCloseModal = () => dispatch('onCloseModal');

    onMount(() => {
        visible = true;
    });

    const closeModal = () => {
        visible = false;
        onCloseModal();
    };

    const handleKeyup = (e) => {
        if (visible && e.key === 'Escape') {
            e.preventDefault();
            closeModal();
        }
    };

    const onClickOutside = () => {
        closeModal();
    };
</script>

{#if visible}
    <div id="modal-wrap"
         class="overlay-modal"
         transition:fade={{duration: 200}}
         on:introstart={onBeforeShowModal}
         on:introend={onAfterShowModal}
         on:outrostart={onBeforeCloseModal}
         on:outroend={onAfterCloseModal}>

        <div class="modal-content {additionalClass}"
             use:clickOutside
             on:click_outside={onClickOutside}
             style={`width: ${width}px`}>

            <div class="modal-header">
                <button type="button" class="btn-close-modal"
                        on:click={closeModal}>
                    <i class="before"></i>
                    <i class="after"></i>
                </button>
            </div>

            <div class="modal-body">
                <slot>
                    <em>no content was provided</em>
                </slot>
            </div>
        </div>

    </div>
{/if}