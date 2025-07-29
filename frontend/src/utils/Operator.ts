/*
 * @Author: harvey_li
 * @Date: 2022-08-03 16:26:23
 * @LastEditors: harvey_li
 * @LastEditTime: 2022-08-05 17:51:21
 * @FilePath: \frontend\src\utils\Operator.ts
 * @Description:
 */
//解决两个操作使用一个tid的问题
export default class Operator {
  tid!: ReturnType<typeof setTimeout>;
  operatorName!: string | undefined;
  constructor(operatorName?: string) {
    this.operatorName = operatorName;
  }
  debounce(callback: ()=>void, delay?: number) {
    this.tid && clearTimeout(this.tid);
    this.tid = setTimeout(
      () => {
        callback();
      },
      delay !== undefined ? delay : 0
    );
  }
}
