import moment from "moment";
import { TimeDesc } from "../models/types";
import { WEEK_ARR } from "../constants";
import { isToday } from "./isToday";

/**
 * 转换日期（包含周几）
 * @param startTime 开始时间
 * @param endTime 结束时间
 * @returns
 */
export const translateDateWithWeek = (
  startTime: TimeDesc,
  endTime: TimeDesc
) => {
  const desc = isToday(
    startTime?.timestamp ? Number(startTime?.timestamp) * 1000 : undefined
  )
    ? "今天"
    : WEEK_ARR[
        moment(
          startTime?.timestamp ? Number(startTime?.timestamp) * 1000 : undefined
        ).isoWeekday() - 1
      ];
  return `${moment(
    startTime?.timestamp ? Number(startTime.timestamp) * 1000 : undefined
  ).format("MM月DD日")} （${desc}）${moment(
    startTime?.timestamp ? Number(startTime.timestamp) * 1000 : undefined
  ).format("HH:mm")}-${moment(
    endTime?.timestamp ? Number(endTime.timestamp) * 1000 : undefined
  ).format("HH:mm")}`;
};

/**
 * 转换日期
 * @param startTime 开始时间
 * @param endTime 结束时间
 * @returns
 */
export const translateDate = (startTime: string, endTime: string) => {
  return startTime && endTime
    ? `${moment(Number(startTime) * 1000).format("YYYY-MM-DD HH:mm")}-${moment(
        Number(endTime) * 1000
      ).format("HH:mm")}`
    : "-";
};
