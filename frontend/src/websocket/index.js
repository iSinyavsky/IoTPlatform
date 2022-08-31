import {DOMAIN} from "../consts/consts";
import {variables, variableValues} from "../store";
import {params} from "../Components/VariableInfo/VariableInfo.svelte";

const websocketUrl = 'ws://'+DOMAIN+ '/ws/connect';

const onopen = () => {
    console.log('WSS connection established.')
};
const onclose = event => {

}
const onmessage = (event, t) => {
    const res = JSON.parse(event.data);
    console.log("message");
    if(res.type === "newValue"){
        variableValues.update(varValues => {
            let value = {
                Current_Value: res.value,
                CreationDate: new Date(res.date/1000)
            }
            if(!varValues[res.varId]) {
                varValues[res.varId] = [];
            }
            varValues[res.varId].push(value)
            console.log(varValues);
            return varValues
        })
    }

}
const onerror = error => {
    console.log('error: ', error.data);
}
let wss = null;
function init() {
    try {
        wss = new WebSocket(websocketUrl);
        wss.onopen = onopen;
        wss.onclose = onclose;
        wss.onmessage = onmessage;
        wss.onerror = onerror;
    } catch(e) {
        console.log(e);
    }
}

export {
    init
};
