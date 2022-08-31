
import { options, BASE_URL } from '../consts/consts';


export function get(url) {
    return fetch(BASE_URL+ url, {
        ...options,
        method: "get",
    }).then(response =>{
        return response.json()
    })
}


export function post(url, body){
    return fetch(BASE_URL + url, {
        ...options,
        method: "post",
        body: JSON.stringify(body)
    }).then((response) =>{
        return response.json()
    })
}


export function postRaw(url, body){
    return fetch(BASE_URL + url, {
        ...options,
        method: "post",
        body: JSON.stringify(body)
    }).then((response) =>{
        return response
    })
}

export function update(url, body){
    return fetch(BASE_URL + url, {
        ...options,
        method: "update",
        body: JSON.stringify(body)
    }).then((response) =>{
        return response.json()
    })
}

export function put(url, body){
    return fetch(BASE_URL + url, {
        ...options,
        method: "PUT",
        body: JSON.stringify(body)
    }).then((response) =>{
        return response.json()
    })
}

export function patch(url, body){
    return fetch(BASE_URL + url, {
        ...options,
        method: "PATCH",
        body: JSON.stringify(body)
    }).then((response) =>{
        return response.json()
    })
}

export function del(url, body){
    return fetch(BASE_URL + url, {
        ...options,
        method: "DELETE",
        body: JSON.stringify(body)
    }).then((response) =>{
        return response.json()
    })
}

export function setBaseUrl(url){
    this.baseUrl = url;
}


