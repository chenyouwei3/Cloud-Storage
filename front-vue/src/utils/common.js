// 生成分享链接
export const makeSharedLink = (route, token) => {
  const baseUrl = window.location.origin
  return `${baseUrl}/shared/${route}?token=${token}`
} 