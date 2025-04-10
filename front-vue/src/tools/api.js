import  request  from './request'; 

export function login(parameter){

    return request({
        url: "/login",
        method: 'post',
        data:parameter
    })
}

export function register(parameter){
    return request({
        url: "/user/add",
        method: 'post',
        data:parameter
    })
}