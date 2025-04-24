import  request  from './request'; 


/*----------------------------------user---------------------------------*/
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

/*----------------------------------dist---------------------------------*/
export function distList(parameter) {
    // 将参数转换为查询字符串
    const queryString = new URLSearchParams(parameter).toString();

    return request({
        url: `/dist/list?${queryString}`, // 在 URL 中拼接查询字符串
        method: 'get', // 使用 GET 请求
        // data: parameter // 不需要传递 data 字段，GET 请求参数已经通过查询字符串传递
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



export function DropdownMenu(parameter) {
    // 将参数转换为查询字符串
    const queryString = new URLSearchParams(parameter).toString();

    return request({
        url: `/dist/dropdownMenu?${queryString}`, // 在 URL 中拼接查询字符串
        method: 'get', // 使用 GET 请求
        // data: parameter // 不需要传递 data 字段，GET 请求参数已经通过查询字符串传递
    });
}

export function distDownload(parameter){
    return request({
        url: "/dist/download",
        method: 'post',
        data: parameter
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

