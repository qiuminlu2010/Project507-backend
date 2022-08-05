import moment from "moment";
export function formatTime(created_data: number) {
	return moment(created_data * 1000).calendar();
}
