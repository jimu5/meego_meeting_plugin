import { create } from "zustand";

interface ICalendarStore {
  workItemId: number;
  workItemTypeKey: string;
  projectId: string;
  updateSign: boolean;
  userId: string;
  init: ({
    workItemId,
    workItemTypeKey,
    projectId,
    userId,
  }: Pick<
    ICalendarStore,
    "projectId" | "workItemId" | "workItemTypeKey" | "userId"
  >) => void;
  reset: () => void;
  switchUpdateSign: (val: boolean) => void;
}
export const useCalendarStore = create<ICalendarStore>((set) => ({
  workItemId: 0,
  workItemTypeKey: "",
  projectId: "",
  userId: "",
  updateSign: false,
  init: ({ workItemId, workItemTypeKey, projectId, userId }) =>
    set((_state) => ({
      workItemId,
      workItemTypeKey,
      projectId,
      userId,
    })),
  reset: () =>
    set((_state) => ({
      workItemId: 0,
      workItemTypeKey: "",
      projectId: "",
      userId: "",
    })),
  switchUpdateSign: (updateSign: boolean) => set((_state) => ({ updateSign })),
}));
