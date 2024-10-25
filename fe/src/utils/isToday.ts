import moment from "moment";

export const isToday = (timestamp) => {
  return (
    moment(timestamp).format("YYYY-MM-DD") === moment().format("YYYY-MM-DD")
  );
};
