const api_domain = "http://127.0.0.1:8080";
export const login = api_domain + "/user/login";
export const register = api_domain + "/user/register";
export const logout = api_domain + "/user/logout";
export const auth_login = api_domain+"/user/auth";
export const create_group = api_domain + "/group/create";

// user/user_detail_id/friends
//    ws://127.0.0.1:3000/_next/webpack-hmr
export const websocketUrl=(user_detail_id) => `"ws://127.0.0.1:8000?user_id=${user_detail_id}"`;

export const friend_list = (user_detail_id) => `${api_domain}/user/${user_detail_id}/friends`;
export const group_list = (user_detail_id) => `${api_domain}/user/${user_detail_id}/groups`;

export const group_detail = (group_id) => `${api_domain}/group/${group_id}/detail`;
export const friend_detail = (friend_id) => `${api_domain}/friend/${friend_id}/detail`;

export const search_group_or_friend = (query) => `${api_domain}/search?query=${query}`;
