import { ResPage, Message, ImageInfo, VideoInfo, Page } from "@/api/interface/index";
import { API } from "@/api/config/servicePort";
import http from "@/api";
import { ContentTypeEnum } from "@/enums/httpEnum";

/**
 * @name 消息管理模块
 */

// * 获取消息列表
export const getMessagesApi = (params: Message.ReqGetParams) => {
	return http.get<Page.Response<Message.MessageInfo>>(API + `/msg/history`, params);
};

// * 获取消息列表
export const getSessionsApi = (params: Message.ReqGetSessionParams) => {
	return http.get<ResPage<Message.SessionInfo>>(API + `/msg/session`, params);
};

// * 更新已读消息
export const updateUnreadMessageApi = (params: { uid: number; session_uid: number }) => {
	return http.post(API + `/msg/read`, params);
};
// * 上传图片
export const uploadImageApi = (params: any) => {
	return http.post<ImageInfo>(API + `/upload/image`, params, { headers: { "Content-Type": ContentTypeEnum.FORM_DATA } });
};

// * 上传视频
export const uploadVideoApi = (params: any) => {
	return http.post<VideoInfo>(API + `/upload/video`, params, { headers: { "Content-Type": ContentTypeEnum.FORM_DATA } });
};
