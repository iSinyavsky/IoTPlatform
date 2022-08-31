import {onDestroy, createEventDispatcher} from 'svelte';

let instanceModal = null;

export function openModal(ConstructorModal, props = {}, callback) {

  props['callback'] = callback;
  return instanceModal = new ConstructorModal({
    target: document.body,
    props,
  });

}