import { writable } from 'svelte/store';
import * as lrFetch from "../common/lrFetch";


export const profile = writable(null);
export const variables = createVariablesStore();
export const variableValues = createVariableValuesStore();

function createVariableValuesStore() {
    const {subscribe, set, update} = writable([]);
    return {
        subscribe,
        set,
        update,
        getVariableValues: (varID) => {
            lrFetch.get("/api/getValues?varid="+varID).then((response) => {
                update(res => {
                    response.map((val)=>{
                        val.CreationDate = new Date(val.CreationDate).getTime()
                    })
                    res[varID] = response
                    return res
                });
            });
        }
    }
}

function createVariablesStore() {
    const {subscribe, set, update} = writable({});
    return {
        subscribe,
        set,
        update,
        getVariables: () => {
            lrFetch.get("/api/getVariables").then((response) => {
                set(response);
            });
        }
    }
}