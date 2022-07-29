/**
 * @param {number} time
 * @param {string} option
 * @returns {string}
 */
import moment from "moment";
export function formatTime(time: any) {
	return moment(time * 1000).calendar();
}
