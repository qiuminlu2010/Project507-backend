import { Article } from "@/api/interface";
import { getArticleListApi, addArticleLike } from "@/api/modules/article";
import { getCommentLike } from "@/api/modules/comment";
import { GlobalStore } from "@/store";
// import axios from "axios";
import { CommentCard, ViewCard } from "./interface";
import { formatTime } from "./utils";
/**
 * 获取随机ID
 * @param {*} length
 * @returns
 */

const globalStore = GlobalStore();
export function randomID(length = 6) {
	return Number(Math.random().toString().substr(3, length) + Date.now()).toString(36);
}

const COLORS = ["#409EFF", "#67C23A", "#E6A23C", "#F56C6C", "#909399"];

function getRandomNum(min: number, max: number) {
	return Math.floor(Math.random() * (max - min + 1)) + min;
}

function randomColor() {
	return COLORS[getRandomNum(0, 4)];
}

export async function getArtileList(offset: number, limit = 1) {
	const list: ViewCard[] = [];
	let params: Article.ReqGetArticleParams = {
		page_num: offset,
		page_size: limit,
		uid: globalStore.uid
	};
	const res = await getArticleListApi(params);
	if (res.code !== 200 || !res.data?.datalist || res.data?.datalist.length == 0) return [];
	res.data?.datalist.forEach(item => {
		if (item.images && item.images.length !== 0) {
			list.push({
				id: item.id,
				star: item.is_like,
				like: item.like_count,
				image_url: `/base` + item.images[0].url,
				src: {
					// original: Math.random() < 0.95 ? successURL : errorURL
					original: `/base` + item.images[0].url
				},
				backgroundColor: randomColor(),
				title: item.title,
				content: item.content,
				data: formatTime(item.created_on * 1000),
				owner_id: item.owner_id,
				owner_name: item.owner_name,
				owner_username: item.owner_username,
				owner_avatar: item.owner_avatar
			});
		}
	});
	return list;
}

// export async function getCommentList(id: number) {
// 	// const end = start + pageSize;
// 	let params = {
// 		article_id: id,
// 		user_id: globalStore.uid
// 	};
// 	const res = await getArticleCommentApi(params);
// 	console.log("res", res);
// 	console.log("id", id);
// 	if (res.data?.datalist.length == 0) {
// 		return null;
// 	} else {
// 		const list: CommentCard[] = [
// 			{
// 				ID: res.data!.datalist[0].ID,
// 				created_on: res.data!.datalist[0].created_on,
// 				user_id: res.data!.datalist[0].user_id,
// 				// article_id: res.data!.datalist[0].article_id,
// 				article_id: id,
// 				reply_id: res.data!.datalist[0]!.reply_id,
// 				username: res.data!.datalist[0].username,
// 				avatar: res.data!.datalist[0].avatar,
// 				content: res.data!.datalist[0].content,
// 				replies: res.data!.datalist[0].replies
// 			}
// 		];
// 		return list;
// 	}
// }

export function handleStar(item: ViewCard) {
	if (item.star) {
		item.like -= 1;
	} else {
		item.like += 1;
	}
	console.log(item.star);
	item.star = !item.star;
	let params = { id: item.id, user_id: globalStore.uid, type: item.star ? 1 : 0 };
	console.log(params, item.star);
	addArticleLike(params);
}

export function handleCommentLike(item: CommentCard) {
	if (item.is_like) {
		item.like! -= 1;
	} else {
		item.like! += 1;
	}
	item.is_like = !item.is_like;
	let params = { id: item.ID, user_id: globalStore.uid, type: item.is_like ? 1 : 0, token: globalStore.token };
	getCommentLike(params);
}
