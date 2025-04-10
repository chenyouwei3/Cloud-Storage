// 用于生成带有提取码的共享链接
export function makeSharedLink(route, token) {
    return `链接: ${route}  提取码: ${token}`;
  }
  