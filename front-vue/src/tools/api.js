import { axios as request } from './request'; 

/*----------------------------------权限中心---------------------------------*/
export function login(parameter){
    return request({
        url: "/login",
        method: 'post',
        data:parameter
    })
}

export function register(parameter){
    return request({
        url: "/user/insert",
        method: 'post',
        data:parameter
    })
}
//权限中心-api
export function  apiList(parameter){
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/api/getList?${queryString}`, 
        method: 'get',
    });
}

export function apiRemove (parameter) {
    return request({
        url: "/api/remove",
        method: 'post',
        data: parameter
    })
}
export function apiInsert(parameter){
    return request({
        url: "/api/insert",
        method: 'post',
        data: parameter
    })
}

export function apiEdit(parameter){
    return request({
        url: "/api/edit",
        method: 'post',
        data: parameter
    })
}
//权限中心-role
export function roleList(parameter){
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/role/getList?${queryString}`, 
        method: 'get', 
    });
}
export function roleRemove (parameter) {
    return request({
        url: "/role/remove",
        method: 'post',
        data: parameter
    })
}
export function roleInsert(parameter){
    return request({
        url: "/role/insert",
        method: 'post',
        data: parameter
    })
}

export function roleEdit(parameter){
    return request({
        url: "/role/edit",
        method: 'post',
        data: parameter
    })
}

export function getRoleByApis(roleId) {
    return request({
        url: `/role/getRoleByApis?id=${roleId}`,
        method: 'get'
    })
}

//权限中心-user
export function userList(parameter){
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/user/getList?${queryString}`, 
        method: 'get', 
    });
}
export function userRemove (parameter) {
    return request({
        url: "/user/remove",
        method: 'post',
        data: parameter
    })
}
export function userInsert(parameter){
    return request({
        url: "/user/insert",
        method: 'post',
        data: parameter
    })
}

export function userEdit(parameter){
    return request({
        url: "/user/edit",
        method: 'post',
        data: parameter
    })
}

export function getUserByRoles(roleId) {
    return request({
        url: `/user/getUserByRoles?id=${roleId}`,
        method: 'get'
    })
}

//日志中心
export function getOperationLog(parameter){
    const queryString = new URLSearchParams(parameter).toString();    // 将参数转换为查询字符串
    return request({
        url: `/log/operation/getList?${queryString}`,
        method: 'get'
    })
}

/*----------------------------------文件云盘---------------------------------*/
export function distList(parameter) {
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/dist/list?${queryString}`, 
        method: 'get',
    });
}

export function distMkdir(parameter){
    return request({
        url: "/dist/mkdir",
        method: 'post',
        data: parameter
    })
}

export function distRename (parameter) {
    return request({
        url: "/dist/rename",
        method: 'post',
        data: parameter
    })
}

export function distRemove (parameter) {
    return request({
        url: "/dist/remove",
        method: 'post',
        data: parameter
    })
}

export function distCopy (parameter) {
    return request({
        url: "/dist/copy",
        method: 'post',
        data: parameter
    })
}

export function distMove (parameter) {
    return request({
        url: "/dist/move",
        method: 'post',
        data: parameter
    })
}

export function distDropdownMenu(parameter) {
    const queryString = new URLSearchParams(parameter).toString();   
    return request({
        url: `/dist/dropdownMenu?${queryString}`, 
        method: 'get', 
    });
}

export function distDownload(parameter){
    return request({
        url: "/dist/download",
        method: 'post',
        data: parameter,
        responseType: 'blob'  // 添加这个配置以正确处理二进制响应
    })
}

export function distUpload(file, path) {
  const formData = new FormData()
  formData.append('file', file)        // 这是上传的文件
  formData.append('path', path)        // 这是文件要保存的路径
  return request({
    url: "/dist/upload",
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function distCategory(parameter) {
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/dist/category?${queryString}`, 
        method: 'get', 
    });
}
