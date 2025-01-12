const api_domain = "http://172.25.59.171:8080";
export const login = api_domain + "/user/login";
export const register = api_domain + "/user/register";
export const logout = api_domain + "/user/logout";
export const auth_login = api_domain+"/user/auth";
export const createGroupUrl = api_domain + "/group/create";

export const applicationUrl = api_domain + "/invite/application";
export const invitionUrl = api_domain + "/invite/invitation";
export const auth_inviteUrl = api_domain + "/invite/auth";
//创建群聊


// user/user_detail_id/friends
//    ws://127.0.0.1:3000/_next/webpack-hmr
export const websocketUrl=(user_detail_id) => `"ws://172.25.59.171:8000?token=${user_detail_id}"`;
export const friend_list = (user_detail_id) => `${api_domain}/user/${user_detail_id}/friends`;
export const group_list = (user_detail_id) => `${api_domain}/user/${user_detail_id}/groups`;

export const group_detail = (group_id) => `${api_domain}/group/${group_id}/detail`;
export const friend_detail = (friend_id) => `${api_domain}/friend/${friend_id}/detail`;

export const search_group_or_friend = (query,type) => `${api_domain}/search/?query=${query}&type=${type}`;
// 删除好友
export const delete_friend = (user_detail_id,friend_id) => `${api_domain}/friend/delete?user_detail_id=${user_detail_id}&friend_id=${friend_id}`;
// 删除群
export const delete_group = (group_id) => `${api_domain}/group/${group_id}/delete`;
// 退出群
export const exit_group = (group_id) => `${api_domain}/group/${group_id}/exit`;



