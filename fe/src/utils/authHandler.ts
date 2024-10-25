export const authHandler = (url?: string) => {
  const aTag = document.createElement("a");
  // 设置a标签的href属性
  aTag.setAttribute("href", url ?? "http://www.example.com");
  // 设置a标签的目标
  aTag.setAttribute("target", "_self");
  aTag.href = "http://www.baidu.com";
  //   window.parent.document.body.appendChild(aTag);
  //   aTag.click();
};
